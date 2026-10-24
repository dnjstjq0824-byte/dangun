import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';
import '../marketplace/listing_detail_screen.dart';

class HomeScreen extends StatefulWidget {
  final ApiClient apiClient;

  const HomeScreen({super.key, required this.apiClient});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  List<Listing> listings = [];
  bool isLoading = true;
  String selectedCategory = '';

  final List<Map<String, String>> categories = [
    {'id': '', 'name': '전체'},
    {'id': 'c0000000-0000-0000-0000-000000000001', 'name': '디지털기기'},
    {'id': 'c0000000-0000-0000-0000-000000000002', 'name': '생활가전'},
    {'id': 'c0000000-0000-0000-0000-000000000003', 'name': '가구/인테리어'},
    {'id': 'c0000000-0000-0000-0000-000000000007', 'name': '도서/음반'},
  ];

  @override
  void initState() {
    super.initState();
    loadListings();
  }

  Future<void> loadListings() async {
    setState(() => isLoading = true);
    try {
      String endpoint = '/listings?sort=latest';
      if (selectedCategory.isNotEmpty) {
        endpoint += '&category_id=$selectedCategory';
      }
      final data = await widget.apiClient.get(endpoint);
      if (data is List) {
        setState(() {
          listings = data.map((e) => Listing.fromJson(e)).toList();
          isLoading = false;
        });
      }
    } catch (e) {
      setState(() => isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        titleSpacing: 16,
        title: Row(
          children: [
            const Text('청라동', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
            const SizedBox(width: 4),
            const Icon(Icons.keyboard_arrow_down, size: 20),
            const SizedBox(width: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
              decoration: BoxDecoration(
                color: const Color(0xFFECFDF5),
                borderRadius: BorderRadius.circular(4),
                border: Border.all(color: const Color(0xFFA7F3D0)),
              ),
              child: const Text('인증됨', style: TextStyle(color: Color(0xFF16A34A), fontSize: 10, fontWeight: FontWeight.w700)),
            ),
          ],
        ),
        actions: [
          IconButton(icon: const Icon(Icons.search), onPressed: () {}),
          IconButton(icon: const Icon(Icons.notifications_none), onPressed: () {}),
        ],
      ),
      body: Column(
        children: [
          // Category Selector
          Container(
            height: 48,
            padding: const EdgeInsets.symmetric(vertical: 6),
            child: ListView.separated(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: categories.length,
              separatorBuilder: (_, _) => const SizedBox(width: 8),
              itemBuilder: (context, index) {
                final cat = categories[index];
                final isSelected = selectedCategory == cat['id'];
                return ChoiceChip(
                  label: Text(cat['name']!),
                  selected: isSelected,
                  selectedColor: AppTheme.brandPrimary,
                  labelStyle: TextStyle(
                    color: isSelected ? Colors.white : AppTheme.textSecondary,
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                  onSelected: (selected) {
                    setState(() => selectedCategory = cat['id']!);
                    loadListings();
                  },
                );
              },
            ),
          ),
          const Divider(height: 1, color: AppTheme.borderDefault),

          // Listings List
          Expanded(
            child: isLoading
                ? const Center(child: CircularProgressIndicator())
                : RefreshIndicator(
                    onRefresh: loadListings,
                    child: listings.isEmpty
                        ? const Center(child: Text('등록된 물품이 없습니다.'))
                        : ListView.separated(
                            itemCount: listings.length,
                            separatorBuilder: (_, _) => const Divider(height: 1, color: AppTheme.borderDefault),
                            itemBuilder: (context, index) {
                              final item = listings[index];
                              return InkWell(
                                onTap: () {
                                  Navigator.push(
                                    context,
                                    MaterialPageRoute(
                                      builder: (_) => ListingDetailScreen(listing: item, apiClient: widget.apiClient),
                                    ),
                                  );
                                },
                                child: Padding(
                                  padding: const EdgeInsets.all(16.0),
                                  child: Row(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      // Thumbnail
                                      ClipRRect(
                                        borderRadius: BorderRadius.circular(8),
                                        child: Image.network(
                                          item.images.isNotEmpty
                                              ? item.images.first
                                              : 'https://images.unsplash.com/photo-1523275335684-37898b6baf30',
                                          width: 100,
                                          height: 100,
                                          fit: BoxFit.cover,
                                          errorBuilder: (_, _, _) => Container(
                                            width: 100,
                                            height: 100,
                                            color: Colors.grey[200],
                                            child: const Icon(Icons.image, color: Colors.grey),
                                          ),
                                        ),
                                      ),
                                      const SizedBox(width: 14),

                                      // Details
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              item.title,
                                              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 15),
                                              maxLines: 2,
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                            const SizedBox(height: 4),
                                            Text(
                                              '${item.neighborhoodName} · ${item.condition}',
                                              style: const TextStyle(color: AppTheme.textMuted, fontSize: 12),
                                            ),
                                            const SizedBox(height: 6),
                                            Text(
                                              item.isFree ? '무료나눔' : '${item.priceMinor.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}원',
                                              style: TextStyle(
                                                color: item.isFree ? AppTheme.statusSuccess : AppTheme.textPrimary,
                                                fontWeight: FontWeight.w800,
                                                fontSize: 16,
                                              ),
                                            ),
                                            const SizedBox(height: 6),
                                            Row(
                                              mainAxisAlignment: MainAxisAlignment.end,
                                              children: [
                                                if (item.chatCount > 0) ...[
                                                  const Icon(Icons.chat_bubble_outline, size: 14, color: AppTheme.textMuted),
                                                  const SizedBox(width: 2),
                                                  Text('${item.chatCount}', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                                                  const SizedBox(width: 8),
                                                ],
                                                if (item.favoriteCount > 0) ...[
                                                  const Icon(Icons.favorite_border, size: 14, color: AppTheme.textMuted),
                                                  const SizedBox(width: 2),
                                                  Text('${item.favoriteCount}', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                                                ],
                                              ],
                                            )
                                          ],
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              );
                            },
                          ),
                  ),
          ),
        ],
      ),
    );
  }
}
