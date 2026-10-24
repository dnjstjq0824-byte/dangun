import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';
import '../chat/chat_room_screen.dart';

class ListingDetailScreen extends StatefulWidget {
  final Listing listing;
  final ApiClient apiClient;

  const ListingDetailScreen({super.key, required this.listing, required this.apiClient});

  @override
  State<ListingDetailScreen> createState() => _ListingDetailScreenState();
}

class _ListingDetailScreenState extends State<ListingDetailScreen> {
  late bool isFavorited;
  late int favCount;

  @override
  void initState() {
    super.initState();
    isFavorited = widget.listing.isFavorited;
    favCount = widget.listing.favoriteCount;
  }

  void toggleFavorite() async {
    setState(() {
      isFavorited = !isFavorited;
      favCount += isFavorited ? 1 : -1;
    });
    try {
      await widget.apiClient.post('/me/favorites', {'listing_id': widget.listing.id});
    } catch (_) {}
  }

  void startChat() async {
    try {
      final res = await widget.apiClient.post('/conversations', {'listing_id': widget.listing.id});
      if (res != null) {
        final conv = Conversation.fromJson(res);
        if (mounted) {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (_) => ChatRoomScreen(conversation: conv, apiClient: widget.apiClient),
            ),
          );
        }
      }
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('채팅 시작 실패: $e')));
    }
  }

  void showOfferDialog() {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('가격 제안하기', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 18)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('희망하는 구매 가격을 입력하세요 (원)', style: TextStyle(fontSize: 13, color: AppTheme.textSecondary)),
            const SizedBox(height: 12),
            TextField(
              controller: controller,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                hintText: '예: 350000',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              final amount = int.tryParse(controller.text) ?? 0;
              Navigator.pop(ctx);
              try {
                await widget.apiClient.post('/listings/${widget.listing.id}/price-offers', {'amount_minor': amount});
                if (!mounted) return;
                ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('가격 제안이 전송되었습니다.')));
              } catch (e) {
                if (!mounted) return;
                ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('제안 실패: $e')));
              }
            },
            child: const Text('제안 전송'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final item = widget.listing;

    return Scaffold(
      appBar: AppBar(
        title: const Text('상품 상세', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        actions: [
          IconButton(icon: const Icon(Icons.share_outlined), onPressed: () {}),
          IconButton(icon: const Icon(Icons.more_vert), onPressed: () {}),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Cover Image
            AspectRatio(
              aspectRatio: 1,
              child: Image.network(
                item.images.isNotEmpty ? item.images.first : 'https://images.unsplash.com/photo-1523275335684-37898b6baf30',
                fit: BoxFit.cover,
                errorBuilder: (_, _, _) => Container(color: Colors.grey[200]),
              ),
            ),

            // Seller Profile Ribbon
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                children: [
                  const CircleAvatar(
                    backgroundColor: AppTheme.surfaceSubtle,
                    child: Icon(Icons.person, color: AppTheme.brandPrimary),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(item.sellerNickname, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15)),
                        Text(item.neighborhoodName, style: const TextStyle(color: AppTheme.textMuted, fontSize: 12)),
                      ],
                    ),
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text('${item.sellerMannerTemp.toStringAsFixed(1)}°C',
                          style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15, color: AppTheme.brandPrimary)),
                      const Text('매너온도', style: TextStyle(color: AppTheme.textMuted, fontSize: 10)),
                    ],
                  ),
                ],
              ),
            ),
            const Divider(height: 1, color: AppTheme.borderDefault),

            // Item Details
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(item.title, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 20)),
                  const SizedBox(height: 8),
                  Text('${item.categoryName} · ${item.condition}',
                      style: const TextStyle(color: AppTheme.textMuted, fontSize: 13)),
                  const SizedBox(height: 16),
                  Text(
                    item.description,
                    style: const TextStyle(fontSize: 15, height: 1.6, color: AppTheme.textPrimary),
                  ),
                  const SizedBox(height: 24),

                  // Safe Area Banner
                  Container(
                    padding: const EdgeInsets.all(14),
                    decoration: BoxDecoration(
                      color: AppTheme.surfaceSubtle,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(color: AppTheme.borderDefault),
                    ),
                    child: Row(
                      children: [
                        const Icon(Icons.shield_outlined, color: AppTheme.brandPrimary, size: 28),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Text('BN NEAR 안심 직거래', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 13)),
                              Text('${item.neighborhoodName} 공공장소(지하철역 등)에서 만나 안전하게 거래하세요.',
                                  style: const TextStyle(fontSize: 11, color: AppTheme.textSecondary)),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
      bottomNavigationBar: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: const BoxDecoration(
          color: Colors.white,
          border: Border(top: BorderSide(color: AppTheme.borderDefault)),
        ),
        child: SafeArea(
          child: Row(
            children: [
              IconButton(
                icon: Icon(isFavorited ? Icons.favorite : Icons.favorite_border,
                    color: isFavorited ? AppTheme.statusDanger : AppTheme.textSecondary),
                onPressed: toggleFavorite,
              ),
              const SizedBox(width: 8),
              Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    item.isFree ? '무료나눔' : '${item.priceMinor.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}원',
                    style: TextStyle(
                      fontWeight: FontWeight.w800,
                      fontSize: 18,
                      color: item.isFree ? AppTheme.statusSuccess : AppTheme.textPrimary,
                    ),
                  ),
                  if (item.isNegotiable && !item.isFree)
                    const Text('가격제안 가능', style: TextStyle(fontSize: 11, color: AppTheme.brandAccent, fontWeight: FontWeight.w600)),
                ],
              ),
              const Spacer(),
              if (item.isNegotiable && !item.isFree) ...[
                OutlinedButton(
                  onPressed: showOfferDialog,
                  style: OutlinedButton.styleFrom(padding: const EdgeInsets.symmetric(horizontal: 12)),
                  child: const Text('가격제안', style: TextStyle(fontSize: 13)),
                ),
                const SizedBox(width: 8),
              ],
              ElevatedButton(
                onPressed: startChat,
                style: ElevatedButton.styleFrom(padding: const EdgeInsets.symmetric(horizontal: 18)),
                child: const Text('채팅하기'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
