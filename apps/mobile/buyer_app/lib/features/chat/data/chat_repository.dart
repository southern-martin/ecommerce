import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

class ChatRepository {
  final ApiClient _apiClient;

  ChatRepository(this._apiClient);

  Future<List<Conversation>> getConversations() async {
    final response = await _apiClient.get('/chat/conversations');
    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => Conversation.fromJson(json)).toList();
  }

  Future<List<Message>> getMessages(String conversationId, {int page = 1}) async {
    final response = await _apiClient.get(
      '/chat/conversations/$conversationId/messages',
      queryParameters: {'page': page},
    );
    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => Message.fromJson(json)).toList();
  }

  Future<Message> sendMessage(String conversationId, String content, {String type = 'text'}) async {
    final response = await _apiClient.post(
      '/chat/conversations/$conversationId/messages',
      data: {'content': content, 'type': type},
    );
    return Message.fromJson(response.data['data']);
  }

  Future<Conversation> createConversation(String participantId, {String type = 'buyer_seller'}) async {
    final response = await _apiClient.post(
      '/chat/conversations',
      data: {'participant_id': participantId, 'type': type},
    );
    return Conversation.fromJson(response.data['data']);
  }
}
