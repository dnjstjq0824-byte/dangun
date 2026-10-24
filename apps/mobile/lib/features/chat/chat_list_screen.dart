import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';
import 'chat_room_screen.dart';

class ChatListScreen extends StatefulWidget {
  final ApiClient apiClient;

  const ChatListScreen({super.key, required this.apiClient});

  @override
  State<ChatListScreen> createState() => _ChatListScreenState();
}

class _ChatListScreenState extends State<ChatListScreen> {
  List<Conversation> conversations = [];
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    loadConversations();
  }

  Future<void> loadConversations() async {
    setState(() => isLoading = true);
    try {
      final data = await widget.apiClient.get('/conversations');
      if (data is List) {
        setState(() {
          conversations = data.map((e) => Conversation.fromJson(e)).toList();
          isLoading = false;
        });
      }
    } catch (_) {
      setState(() => isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('채팅', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: loadConversations,
              child: conversations.isEmpty
                  ? const Center(child: Text('진행 중인 대화가 없습니다.'))
                  : ListView.separated(
                      itemCount: conversations.length,
                      separatorBuilder: (_, _) => const Divider(height: 1, color: AppTheme.borderDefault),
                      itemBuilder: (context, index) {
                        final conv = conversations[index];
                        return ListTile(
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                          leading: const CircleAvatar(
                            backgroundColor: AppTheme.surfaceSubtle,
                            child: Icon(Icons.person, color: AppTheme.brandPrimary),
                          ),
                          title: Row(
                            children: [
                              Text(conv.sellerNickname, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15)),
                              const SizedBox(width: 6),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
                                decoration: BoxDecoration(
                                  color: AppTheme.surfaceSubtle,
                                  borderRadius: BorderRadius.circular(4),
                                ),
                                child: Text(conv.listingStatus, style: const TextStyle(fontSize: 10, fontWeight: FontWeight.w600)),
                              ),
                            ],
                          ),
                          subtitle: Text(
                            conv.lastMessage,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(fontSize: 13, color: AppTheme.textSecondary),
                          ),
                          trailing: ClipRRect(
                            borderRadius: BorderRadius.circular(6),
                            child: Image.network(
                              conv.listingImage.isNotEmpty ? conv.listingImage : 'https://images.unsplash.com/photo-1523275335684-37898b6baf30',
                              width: 48,
                              height: 48,
                              fit: BoxFit.cover,
                              errorBuilder: (_, _, _) => Container(width: 48, height: 48, color: Colors.grey[200]),
                            ),
                          ),
                          onTap: () async {
                            await Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (_) => ChatRoomScreen(conversation: conv, apiClient: widget.apiClient),
                              ),
                            );
                            loadConversations();
                          },
                        );
                      },
                    ),
            ),
    );
  }
}
