import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';

class ChatRoomScreen extends StatefulWidget {
  final Conversation conversation;
  final ApiClient apiClient;

  const ChatRoomScreen({super.key, required this.conversation, required this.apiClient});

  @override
  State<ChatRoomScreen> createState() => _ChatRoomScreenState();
}

class _ChatRoomScreenState extends State<ChatRoomScreen> {
  final List<ChatMessage> messages = [];
  final TextEditingController _msgController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  WebSocketChannel? channel;
  String dealStatus = 'ACTIVE';

  @override
  void initState() {
    super.initState();
    dealStatus = widget.conversation.listingStatus;
    loadMessages();
    connectWebSocket();
  }

  @override
  void dispose() {
    channel?.sink.close();
    _msgController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void connectWebSocket() {
    final token = widget.apiClient.token;
    if (token == null || token.isEmpty) return;

    try {
      final wsUri = Uri.parse('${ApiClient.wsUrl}?token=$token');
      channel = WebSocketChannel.connect(wsUri);
      channel!.stream.listen((event) {
        try {
          final data = jsonDecode(event);
          if (data['type'] == 'chat.message.receive' || data['type'] == 'chat.message.ack') {
            if (data['conversation_id'] == widget.conversation.id) {
              final newSeq = data['sequence_no'] as int;
              if (!messages.any((m) => m.sequenceNo == newSeq)) {
                setState(() {
                  messages.add(ChatMessage(
                    id: data['message_id'] ?? '',
                    conversationId: widget.conversation.id,
                    senderUserId: data['sender_user_id'] ?? '',
                    sequenceNo: newSeq,
                    messageType: 'TEXT',
                    body: data['payload']?['text'] ?? '',
                    createdAt: DateTime.now(),
                  ));
                });
                scrollToBottom();
              }
            }
          }
        } catch (_) {}
      });
    } catch (_) {}
  }

  Future<void> loadMessages() async {
    try {
      final data = await widget.apiClient.get('/conversations/${widget.conversation.id}/messages');
      if (data is List) {
        setState(() {
          messages.clear();
          messages.addAll(data.map((e) => ChatMessage.fromJson(e)));
        });
        scrollToBottom();
      }
    } catch (_) {}
  }

  void scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void sendMessage() async {
    final text = _msgController.text.trim();
    if (text.isEmpty) return;
    _msgController.clear();

    if (channel != null) {
      final payload = {
        'type': 'chat.message.send',
        'conversation_id': widget.conversation.id,
        'client_message_id': DateTime.now().millisecondsSinceEpoch.toString(),
        'payload': {
          'message_type': 'TEXT',
          'text': text,
        }
      };
      channel!.sink.add(jsonEncode(payload));
    } else {
      // Fallback REST
      try {
        final res = await widget.apiClient.post(
          '/conversations/${widget.conversation.id}/messages',
          {
            'client_message_id': DateTime.now().millisecondsSinceEpoch.toString(),
            'message_type': 'TEXT',
            'body': text,
          },
        );
        if (res != null) {
          setState(() => messages.add(ChatMessage.fromJson(res)));
          scrollToBottom();
        }
      } catch (_) {}
    }
  }

  void reserveItem() async {
    try {
      await widget.apiClient.post('/deals/${widget.conversation.listingId}/reserve', {
        'buyer_user_id': widget.conversation.buyerUserId,
      });
      if (mounted) {
        setState(() => dealStatus = 'RESERVED');
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('거래 예약되었습니다.')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('예약 실패: $e')));
      }
    }
  }

  void markSold() async {
    try {
      await widget.apiClient.post('/deals/${widget.conversation.listingId}/seller-complete', {});
      if (mounted) {
        setState(() => dealStatus = 'COMPLETED_CLAIMED');
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('판매 완료되었습니다.')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('완료 처리 실패: $e')));
      }
    }
  }

  void confirmPurchase() async {
    try {
      await widget.apiClient.post('/deals/${widget.conversation.listingId}/buyer-confirm', {});
      if (mounted) {
        setState(() => dealStatus = 'COMPLETED_CONFIRMED');
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('구매 확정되었습니다.')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('확정 실패: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.conversation.sellerNickname, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        actions: [
          IconButton(icon: const Icon(Icons.more_vert), onPressed: () {}),
        ],
      ),
      body: Column(
        children: [
          // Top Listing Header
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
            color: AppTheme.surfaceSubtle,
            child: Row(
              children: [
                ClipRRect(
                  borderRadius: BorderRadius.circular(6),
                  child: Image.network(
                    widget.conversation.listingImage.isNotEmpty ? widget.conversation.listingImage : 'https://images.unsplash.com/photo-1523275335684-37898b6baf30',
                    width: 44,
                    height: 44,
                    fit: BoxFit.cover,
                    errorBuilder: (_, _, _) => Container(width: 44, height: 44, color: Colors.grey[200]),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(widget.conversation.listingTitle, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 13), maxLines: 1),
                      Text('${widget.conversation.listingPrice}원', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 14, color: AppTheme.brandPrimary)),
                    ],
                  ),
                ),
                // Deal Action Button
                if (dealStatus == 'ACTIVE')
                  ElevatedButton(
                    onPressed: reserveItem,
                    style: ElevatedButton.styleFrom(padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6)),
                    child: const Text('예약하기', style: TextStyle(fontSize: 12)),
                  )
                else if (dealStatus == 'RESERVED')
                  ElevatedButton(
                    onPressed: markSold,
                    style: ElevatedButton.styleFrom(backgroundColor: AppTheme.statusReserved, padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6)),
                    child: const Text('판매완료', style: TextStyle(fontSize: 12)),
                  )
                else if (dealStatus == 'COMPLETED_CLAIMED')
                  ElevatedButton(
                    onPressed: confirmPurchase,
                    style: ElevatedButton.styleFrom(backgroundColor: AppTheme.statusSuccess, padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6)),
                    child: const Text('구매확정', style: TextStyle(fontSize: 12)),
                  )
                else
                  const Chip(label: Text('거래완료', style: TextStyle(fontSize: 11))),
              ],
            ),
          ),
          const Divider(height: 1, color: AppTheme.borderDefault),

          // Message Stream
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: messages.length,
              itemBuilder: (context, index) {
                final msg = messages[index];
                // In demo, check if sender is me or other
                final isMe = msg.senderUserId.contains('buyer') || msg.senderUserId == widget.conversation.buyerUserId;

                return Align(
                  alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.only(bottom: 8),
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                    decoration: BoxDecoration(
                      color: isMe ? AppTheme.brandPrimary : AppTheme.surfaceCard,
                      borderRadius: BorderRadius.circular(14),
                      border: Border.all(color: isMe ? AppTheme.brandPrimary : AppTheme.borderDefault),
                    ),
                    child: Text(
                      msg.body,
                      style: TextStyle(
                        color: isMe ? Colors.white : AppTheme.textPrimary,
                        fontSize: 14,
                      ),
                    ),
                  ),
                );
              },
            ),
          ),

          // Composer
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: const BoxDecoration(
              color: Colors.white,
              border: Border(top: BorderSide(color: AppTheme.borderDefault)),
            ),
            child: SafeArea(
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.add, color: AppTheme.textSecondary),
                    onPressed: () {},
                  ),
                  Expanded(
                    child: TextField(
                      controller: _msgController,
                      decoration: const InputDecoration(
                        hintText: '메시지를 입력하세요...',
                        border: InputBorder.none,
                      ),
                      onSubmitted: (_) => sendMessage(),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.send, color: AppTheme.brandPrimary),
                    onPressed: sendMessage,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
