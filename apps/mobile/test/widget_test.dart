import 'package:flutter_test/flutter_test.dart';
import 'package:bn_near/main.dart';

void main() {
  testWidgets('BN NEAR App smoke and shell render test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(const BnNearApp());

    // Verify Bottom Navigation items exist
    expect(find.text('홈'), findsOneWidget);
    expect(find.text('동네생활'), findsOneWidget);
    expect(find.text('글쓰기'), findsOneWidget);
    expect(find.text('채팅'), findsOneWidget);
    expect(find.text('나의 NEAR'), findsOneWidget);
  });
}
