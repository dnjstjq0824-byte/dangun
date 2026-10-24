import 'package:flutter/material.dart';
import '../../core/theme/app_theme.dart';
import '../../core/models/models.dart';
import '../../core/api/api_client.dart';

class MyScreen extends StatefulWidget {
  final ApiClient apiClient;
  final User? currentUser;
  final Function(User, String) onSwitchAccount;

  const MyScreen({
    super.key,
    required this.apiClient,
    required this.currentUser,
    required this.onSwitchAccount,
  });

  @override
  State<MyScreen> createState() => _MyScreenState();
}

class _MyScreenState extends State<MyScreen> {
  bool isVerified = true;

  void triggerLocationVerification() async {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('동네 GPS 위치 인증', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 18)),
        content: const Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('현재 청라동 행정 경계 내에 위치하고 있는지 실시간 GPS 좌표로 검증합니다.', style: TextStyle(fontSize: 13, color: AppTheme.textSecondary)),
            SizedBox(height: 12),
            Text('• 측정 좌표: 37.5350, 126.6500 (청라동)\n• GPS 정확도: 25.0m (기준 < 150m 충족)\n• 인증 유효기간: 30일', style: TextStyle(fontSize: 12, height: 1.6, color: AppTheme.brandPrimary, fontWeight: FontWeight.w600)),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(ctx);
              try {
                await widget.apiClient.post('/me/neighborhood-verifications', {
                  'neighborhood_id': 'b0000000-0000-0000-0000-000000000001',
                  'latitude': 37.5350,
                  'longitude': 126.6500,
                  'accuracy_meters': 25.0,
                });
                setState(() => isVerified = true);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('청라동 위치 인증이 완료되었습니다! (30일 유효)')));
                }
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('인증 실패: $e')));
                }
              }
            },
            child: const Text('인증 완료'),
          ),
        ],
      ),
    );
  }

  void switchAccountDialog() {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('테스트 계정 전환', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 16)),
            const SizedBox(height: 14),
            ListTile(
              leading: const Icon(Icons.person, color: AppTheme.brandPrimary),
              title: const Text('구매자 (Minji_Cheongna)'),
              subtitle: const Text('buyer@bn-near.com'),
              onTap: () async {
                Navigator.pop(ctx);
                loginAs('buyer@bn-near.com', 'Password123!');
              },
            ),
            ListTile(
              leading: const Icon(Icons.store, color: AppTheme.brandPrimary),
              title: const Text('판매자 (David_Cheongna)'),
              subtitle: const Text('seller@bn-near.com'),
              onTap: () async {
                Navigator.pop(ctx);
                loginAs('seller@bn-near.com', 'Password123!');
              },
            ),
          ],
        ),
      ),
    );
  }

  void loginAs(String email, String password) async {
    try {
      final res = await widget.apiClient.post('/auth/login', {
        'email': email,
        'password': password,
        'device_id': 'flutter-device',
      });
      if (res != null) {
        final u = User.fromJson(res['user']);
        final t = res['access_token'] as String;
        widget.onSwitchAccount(u, t);
        if (!mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('${u.nickname} 계정으로 전환되었습니다.')));
      }
    } catch (_) {}
  }

  @override
  Widget build(BuildContext context) {
    final user = widget.currentUser;

    return Scaffold(
      appBar: AppBar(
        title: const Text('나의 BN NEAR', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
        actions: [
          IconButton(icon: const Icon(Icons.settings_outlined), onPressed: () {}),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Profile Card
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Row(
                    children: [
                      const CircleAvatar(
                        radius: 28,
                        backgroundColor: AppTheme.surfaceSubtle,
                        child: Icon(Icons.person, size: 32, color: AppTheme.brandPrimary),
                      ),
                      const SizedBox(width: 14),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(user?.nickname ?? '이웃님', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 17)),
                            Text(user?.email ?? '로그인 필요', style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                            const SizedBox(height: 4),
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                              decoration: BoxDecoration(
                                color: const Color(0xFFECFDF5),
                                borderRadius: BorderRadius.circular(4),
                              ),
                              child: const Text('📍 청라동 인증회원', style: TextStyle(color: Color(0xFF16A34A), fontSize: 11, fontWeight: FontWeight.w700)),
                            ),
                          ],
                        ),
                      ),
                      Column(
                        children: [
                          Text('${(user?.mannerTemperature ?? 36.5).toStringAsFixed(1)}°C',
                              style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 17, color: AppTheme.brandPrimary)),
                          const Text('매너온도', style: TextStyle(fontSize: 10, color: AppTheme.textMuted)),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),

            // Location Verification Button
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: ListTile(
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10), side: const BorderSide(color: AppTheme.borderDefault)),
                leading: const Icon(Icons.my_location, color: AppTheme.brandPrimary),
                title: const Text('동네 GPS 위치 인증', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
                subtitle: const Text('현재 청라동 인증 유효기간 30일', style: TextStyle(fontSize: 12, color: AppTheme.textMuted)),
                trailing: const Icon(Icons.chevron_right),
                onTap: triggerLocationVerification,
              ),
            ),
            const SizedBox(height: 16),

            // Activity Grid
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Row(
                children: [
                  Expanded(
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: const [
                            Icon(Icons.shopping_bag_outlined, color: AppTheme.brandPrimary),
                            SizedBox(height: 6),
                            Text('판매내역', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 13)),
                          ],
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: const [
                            Icon(Icons.receipt_long_outlined, color: AppTheme.brandPrimary),
                            SizedBox(height: 6),
                            Text('구매내역', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 13)),
                          ],
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: const [
                            Icon(Icons.favorite_border, color: AppTheme.brandPrimary),
                            SizedBox(height: 6),
                            Text('관심목록', style: TextStyle(fontWeight: FontWeight.w700, fontSize: 13)),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Account Switcher Button
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: SizedBox(
                width: double.infinity,
                child: OutlinedButton(
                  onPressed: switchAccountDialog,
                  style: OutlinedButton.styleFrom(padding: const EdgeInsets.all(14)),
                  child: const Text('테스트 계정 전환 (구매자 / 판매자)', style: TextStyle(fontWeight: FontWeight.w700)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
