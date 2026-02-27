import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';

import '../di/injection.dart';

class PushNotificationService {
  static final FlutterLocalNotificationsPlugin _localNotifications =
      FlutterLocalNotificationsPlugin();
  static final FirebaseMessaging _messaging = FirebaseMessaging.instance;
  static GlobalKey<NavigatorState>? _navigatorKey;

  static Future<void> initialize({GlobalKey<NavigatorState>? navigatorKey}) async {
    _navigatorKey = navigatorKey;

    await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings();
    const initSettings = InitializationSettings(
      android: androidSettings,
      iOS: iosSettings,
    );

    await _localNotifications.initialize(
      initSettings,
      onDidReceiveNotificationResponse: _onNotificationTapped,
    );

    const channel = AndroidNotificationChannel(
      'seller_channel',
      'Seller Notifications',
      description: 'Notifications for seller activities',
      importance: Importance.high,
    );
    await _localNotifications
        .resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(channel);

    final token = await _messaging.getToken();
    if (token != null) {
      await _registerToken(token);
    }

    _messaging.onTokenRefresh.listen(_registerToken);
    FirebaseMessaging.onMessage.listen(_showLocalNotification);
    FirebaseMessaging.onMessageOpenedApp.listen(_handleMessageTap);

    final initialMessage = await _messaging.getInitialMessage();
    if (initialMessage != null) {
      _handleMessageTap(initialMessage);
    }
  }

  static Future<void> _registerToken(String token) async {
    try {
      final apiClient = getIt<ApiClient>();
      await apiClient.post('/users/me/fcm-token', data: {'token': token});
    } catch (_) {}
  }

  static void _showLocalNotification(RemoteMessage message) {
    final notification = message.notification;
    if (notification == null) return;

    _localNotifications.show(
      notification.hashCode,
      notification.title,
      notification.body,
      const NotificationDetails(
        android: AndroidNotificationDetails(
          'seller_channel',
          'Seller Notifications',
          importance: Importance.high,
          priority: Priority.high,
          icon: '@mipmap/ic_launcher',
        ),
        iOS: DarwinNotificationDetails(),
      ),
      payload: message.data['type'] ?? '',
    );
  }

  static void _onNotificationTapped(NotificationResponse response) {
    final type = response.payload;
    if (type == null || type.isEmpty) return;
    _navigateByType(type, null);
  }

  static void _handleMessageTap(RemoteMessage message) {
    final type = message.data['type'] as String?;
    final id = message.data['id'] as String?;
    if (type == null) return;
    _navigateByType(type, id);
  }

  static void _navigateByType(String type, String? id) {
    final context = _navigatorKey?.currentContext;
    if (context == null) return;

    switch (type) {
      case 'order':
        if (id != null) context.push('/orders/$id');
        break;
      case 'return':
        if (id != null) context.push('/returns/$id');
        break;
      case 'payout':
        context.push('/payouts');
        break;
      case 'product':
        context.push('/products');
        break;
      default:
        context.push('/dashboard');
    }
  }
}
