import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class AIRepository {
  final ApiClient _apiClient;

  AIRepository(this._apiClient);

  Future<AIResponse> chat(String message, {String? conversationId}) async {
    final response = await _apiClient.post(
      '/ai/chat',
      data: {
        'message': message,
        if (conversationId != null) 'conversation_id': conversationId,
      },
    );
    return AIResponse.fromJson(response.data['data']);
  }

  Future<List<AIConversation>> getConversations() async {
    final response = await _apiClient.get('/ai/conversations');
    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => AIConversation.fromJson(json)).toList();
  }
}

class AIResponse {
  final String conversationId;
  final String reply;
  final List<Map<String, dynamic>>? productRecommendations;

  AIResponse({
    required this.conversationId,
    required this.reply,
    this.productRecommendations,
  });

  factory AIResponse.fromJson(Map<String, dynamic> json) {
    return AIResponse(
      conversationId: json['conversation_id'] ?? '',
      reply: json['reply'] ?? '',
      productRecommendations: json['product_recommendations'] != null
          ? List<Map<String, dynamic>>.from(json['product_recommendations'])
          : null,
    );
  }
}

class AIConversation {
  final String id;
  final String lastMessage;
  final DateTime createdAt;

  AIConversation({required this.id, required this.lastMessage, required this.createdAt});

  factory AIConversation.fromJson(Map<String, dynamic> json) {
    return AIConversation(
      id: json['id'] ?? '',
      lastMessage: json['last_message'] ?? '',
      createdAt: DateTime.parse(json['created_at'] ?? DateTime.now().toIso8601String()),
    );
  }
}
