import 'package:flutter/material.dart';
import 'core/theme/app_theme.dart';
import 'core/api/api_client.dart';
import 'core/models/models.dart';
import 'features/home/home_screen.dart';
import 'features/community/community_screen.dart';
import 'features/sell/sell_screen.dart';
import 'features/chat/chat_list_screen.dart';
import 'features/my/my_screen.dart';

void main() {
  runApp(const BnNearApp());
}

class BnNearApp extends StatelessWidget {
  const BnNearApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'BN NEAR',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: const MainNavigationShell(),
    );
  }
}

class MainNavigationShell extends StatefulWidget {
  const MainNavigationShell({super.key});

  @override
  State<MainNavigationShell> createState() => _MainNavigationShellState();
}

class _MainNavigationShellState extends State<MainNavigationShell> {
  int _currentIndex = 0;
  final ApiClient _apiClient = ApiClient();
  User? _currentUser;

  @override
  void initState() {
    super.initState();
    _initDemoAuth();
  }

  void _initDemoAuth() async {
    try {
      final res = await _apiClient.post('/auth/login', {
        'email': 'buyer@bn-near.com',
        'password': 'Password123!',
        'device_id': 'flutter-app',
      });
      if (res != null) {
        setState(() {
          _apiClient.token = res['access_token'];
          _currentUser = User.fromJson(res['user']);
        });
      }
    } catch (_) {}
  }

  void _handleSwitchAccount(User newUser, String newToken) {
    setState(() {
      _currentUser = newUser;
      _apiClient.token = newToken;
    });
  }

  @override
  Widget build(BuildContext context) {
    final screens = [
      HomeScreen(apiClient: _apiClient),
      CommunityScreen(apiClient: _apiClient),
      SellScreen(apiClient: _apiClient),
      ChatListScreen(apiClient: _apiClient),
      MyScreen(
        apiClient: _apiClient,
        currentUser: _currentUser,
        onSwitchAccount: _handleSwitchAccount,
      ),
    ];

    return Scaffold(
      body: IndexedStack(
        index: _currentIndex,
        children: screens,
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _currentIndex,
        onTap: (index) => setState(() => _currentIndex = index),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home_outlined), activeIcon: Icon(Icons.home), label: '홈'),
          BottomNavigationBarItem(icon: Icon(Icons.forum_outlined), activeIcon: Icon(Icons.forum), label: '동네생활'),
          BottomNavigationBarItem(icon: Icon(Icons.add_circle_outline), activeIcon: Icon(Icons.add_circle), label: '글쓰기'),
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble_outline), activeIcon: Icon(Icons.chat_bubble), label: '채팅'),
          BottomNavigationBarItem(icon: Icon(Icons.person_outline), activeIcon: Icon(Icons.person), label: '나의 NEAR'),
        ],
      ),
    );
  }
}
