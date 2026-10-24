import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

class ApiClient {
  static String get baseUrl {
    if (kIsWeb) {
      return 'http://localhost:8080/api/v1';
    }
    // Default to localhost for desktop and 10.0.2.2 for android
    if (defaultTargetPlatform == TargetPlatform.android) {
      return 'http://10.0.2.2:8080/api/v1';
    }
    return 'http://localhost:8080/api/v1';
  }

  static String get wsUrl {
    if (kIsWeb) {
      return 'ws://localhost:8080/api/v1/ws';
    }
    if (defaultTargetPlatform == TargetPlatform.android) {
      return 'ws://10.0.2.2:8080/api/v1/ws';
    }
    return 'ws://localhost:8080/api/v1/ws';
  }

  String? token;

  ApiClient({this.token});

  Map<String, String> get headers {
    final h = {'Content-Type': 'application/json'};
    if (token != null && token!.isNotEmpty) {
      h['Authorization'] = 'Bearer $token';
    }
    return h;
  }

  Future<dynamic> get(String endpoint) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final response = await http.get(uri, headers: headers);
    return _handleResponse(response);
  }

  Future<dynamic> post(String endpoint, Map<String, dynamic> body) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final response = await http.post(uri, headers: headers, body: jsonEncode(body));
    return _handleResponse(response);
  }

  dynamic _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final decoded = jsonDecode(utf8.decode(response.bodyBytes));
      return decoded['data'] ?? decoded;
    } else {
      throw Exception('API Error: ${response.statusCode} - ${response.body}');
    }
  }
}

