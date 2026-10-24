import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';

class CommunityScreen extends StatefulWidget {
  final ApiClient apiClient;

  const CommunityScreen({super.key, required this.apiClient});

  @override
  State<CommunityScreen> createState() => _CommunityScreenState();
}

class _CommunityScreenState extends State<CommunityScreen> {
  List<CommunityPost> posts = [];
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    loadPosts();
  }

  Future<void> loadPosts() async {
    setState(() => isLoading = true);
    try {
      final data = await widget.apiClient.get('/community/posts?neighborhood_id=b0000000-0000-0000-0000-000000000001');
      if (data is List) {
        setState(() {
          posts = data.map((e) => CommunityPost.fromJson(e)).toList();
          isLoading = false;
        });
      }
    } catch (_) {
      setState(() => isLoading = false);
    }
  }

  void showWriteDialog() {
    final titleCtrl = TextEditingController();
    final bodyCtrl = TextEditingController();
    String category = '동네질문';

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setModalState) => AlertDialog(
          title: const Text('동네생활 글쓰기', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 18)),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              DropdownButton<String>(
                value: category,
                isExpanded: true,
                items: ['동네질문', '생활정보', '동네맛집/추천', '일상'].map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
                onChanged: (val) {
                  if (val != null) setModalState(() => category = val);
                },
              ),
              const SizedBox(height: 8),
              TextField(
                controller: titleCtrl,
                decoration: const InputDecoration(hintText: '제목을 입력하세요'),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: bodyCtrl,
                maxLines: 4,
                decoration: const InputDecoration(hintText: '이웃들과 나누고 싶은 이야기를 적어보세요.'),
              ),
            ],
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
            ElevatedButton(
              onPressed: () async {
                Navigator.pop(ctx);
                try {
                  await widget.apiClient.post('/community/posts', {
                    'neighborhood_id': 'b0000000-0000-0000-0000-000000000001',
                    'category': category,
                    'title': titleCtrl.text,
                    'body': bodyCtrl.text,
                  });
                  loadPosts();
                } catch (_) {}
              },
              child: const Text('등록'),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('동네생활', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
        actions: [
          IconButton(icon: const Icon(Icons.search), onPressed: () {}),
        ],
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: loadPosts,
              child: ListView.separated(
                itemCount: posts.length,
                separatorBuilder: (_, _) => const Divider(height: 1, color: AppTheme.borderDefault),
                itemBuilder: (context, index) {
                  final post = posts[index];
                  return Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                          decoration: BoxDecoration(
                            color: AppTheme.surfaceSubtle,
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(post.category, style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: AppTheme.brandPrimary)),
                        ),
                        const SizedBox(height: 8),
                        Text(post.title, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16)),
                        const SizedBox(height: 6),
                        Text(post.body, style: const TextStyle(fontSize: 14, color: AppTheme.textSecondary, height: 1.5), maxLines: 3, overflow: TextOverflow.ellipsis),
                        const SizedBox(height: 12),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text('${post.authorNickname} · ${post.neighborhoodName}', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                            Row(
                              children: [
                                const Icon(Icons.favorite_border, size: 14, color: AppTheme.textMuted),
                                const SizedBox(width: 4),
                                Text('${post.likeCount}', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                                const SizedBox(width: 12),
                                const Icon(Icons.chat_bubble_outline, size: 14, color: AppTheme.textMuted),
                                const SizedBox(width: 4),
                                Text('${post.commentCount}', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                              ],
                            ),
                          ],
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: showWriteDialog,
        backgroundColor: AppTheme.brandPrimary,
        icon: const Icon(Icons.edit, color: Colors.white),
        label: const Text('글쓰기', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w700)),
      ),
    );
  }
}
