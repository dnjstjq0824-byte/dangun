import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/api/api_client.dart';

class SellScreen extends StatefulWidget {
  final ApiClient apiClient;

  const SellScreen({super.key, required this.apiClient});

  @override
  State<SellScreen> createState() => _SellScreenState();
}

class _SellScreenState extends State<SellScreen> {
  final _titleController = TextEditingController();
  final _priceController = TextEditingController();
  final _descController = TextEditingController();
  final _imgController = TextEditingController(text: 'https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80');

  String categoryId = 'c0000000-0000-0000-0000-000000000001';
  String condition = 'GOOD';
  bool isFree = false;
  bool isNegotiable = true;
  bool isSubmitting = false;

  void submitListing() async {
    if (_titleController.text.isEmpty || _descController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('제목과 설명을 입력해주세요.')));
      return;
    }

    setState(() => isSubmitting = true);
    final price = isFree ? 0 : int.tryParse(_priceController.text) ?? 0;

    try {
      await widget.apiClient.post('/listings', {
        'title': _titleController.text,
        'description': _descController.text,
        'category_id': categoryId,
        'condition': condition,
        'price_minor': price,
        'currency_code': 'KRW',
        'is_free': isFree,
        'is_negotiable': isNegotiable,
        'neighborhood_id': 'b0000000-0000-0000-0000-000000000001',
        'images': [_imgController.text],
      });

      setState(() => isSubmitting = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('물품이 등록되었습니다!')));
        _titleController.clear();
        _priceController.clear();
        _descController.clear();
      }
    } catch (e) {
      setState(() => isSubmitting = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('등록 실패: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('내 물건 판매하기', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Image URL Field
            const Text('대표 사진 URL', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
            const SizedBox(height: 6),
            TextField(
              controller: _imgController,
              decoration: const InputDecoration(border: OutlineInputBorder(), hintText: 'https://...'),
            ),
            const SizedBox(height: 16),

            // Title
            const Text('글 제목', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
            const SizedBox(height: 6),
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(border: OutlineInputBorder(), hintText: '상품명을 입력해주세요'),
            ),
            const SizedBox(height: 16),

            // Category
            const Text('카테고리', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
            const SizedBox(height: 6),
            DropdownButtonFormField<String>(
              initialValue: categoryId,
              decoration: const InputDecoration(border: OutlineInputBorder()),
              items: const [
                DropdownMenuItem(value: 'c0000000-0000-0000-0000-000000000001', child: Text('디지털기기')),
                DropdownMenuItem(value: 'c0000000-0000-0000-0000-000000000002', child: Text('생활가전')),
                DropdownMenuItem(value: 'c0000000-0000-0000-0000-000000000003', child: Text('가구/인테리어')),
                DropdownMenuItem(value: 'c0000000-0000-0000-0000-000000000007', child: Text('도서/음반')),
              ],
              onChanged: (v) {
                if (v != null) setState(() => categoryId = v);
              },
            ),
            const SizedBox(height: 16),

            // Condition
            const Text('상품 상태', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
            const SizedBox(height: 6),
            Wrap(
              spacing: 8,
              children: ['NEW', 'LIKE_NEW', 'GOOD', 'FAIR'].map((c) {
                final isSelected = condition == c;
                return ChoiceChip(
                  label: Text(c),
                  selected: isSelected,
                  selectedColor: AppTheme.brandPrimary,
                  labelStyle: TextStyle(color: isSelected ? Colors.white : AppTheme.textSecondary),
                  onSelected: (s) => setState(() => condition = c),
                );
              }).toList(),
            ),
            const SizedBox(height: 16),

            // Price & Free toggle
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('판매 가격 (원)', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
                Row(
                  children: [
                    Checkbox(
                      value: isFree,
                      activeColor: AppTheme.brandPrimary,
                      onChanged: (v) => setState(() => isFree = v ?? false),
                    ),
                    const Text('무료나눔', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600)),
                  ],
                ),
              ],
            ),
            if (!isFree) ...[
              TextField(
                controller: _priceController,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(border: OutlineInputBorder(), hintText: '₩ 가격 입력'),
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Checkbox(
                    value: isNegotiable,
                    activeColor: AppTheme.brandPrimary,
                    onChanged: (v) => setState(() => isNegotiable = v ?? true),
                  ),
                  const Text('가격 제안 받기', style: TextStyle(fontSize: 13)),
                ],
              ),
            ],
            const SizedBox(height: 16),

            // Description
            const Text('자세한 설명', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
            const SizedBox(height: 6),
            TextField(
              controller: _descController,
              maxLines: 5,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                hintText: '구매 시기, 사용감, 하자 유무 등을 솔직하게 작성해주세요.',
              ),
            ),
            const SizedBox(height: 24),

            // Submit Button
            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: isSubmitting ? null : submitListing,
                child: isSubmitting
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('등록 완료하기', style: TextStyle(fontSize: 16)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
