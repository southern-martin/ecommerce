import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class UserProfile {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final String? avatarUrl;
  final String? phone;
  final DateTime? dateOfBirth;

  const UserProfile({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    this.avatarUrl,
    this.phone,
    this.dateOfBirth,
  });

  String get fullName => '$firstName $lastName';

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      id: json['id'] as String,
      email: json['email'] as String,
      firstName: json['firstName'] as String,
      lastName: json['lastName'] as String,
      avatarUrl: json['avatarUrl'] as String?,
      phone: json['phone'] as String?,
      dateOfBirth: json['dateOfBirth'] != null
          ? DateTime.parse(json['dateOfBirth'] as String)
          : null,
    );
  }
}

class Address {
  final String id;
  final String label;
  final String name;
  final String street;
  final String city;
  final String state;
  final String zip;
  final String country;
  final String? phone;
  final bool isDefault;

  const Address({
    required this.id,
    required this.label,
    required this.name,
    required this.street,
    required this.city,
    required this.state,
    required this.zip,
    required this.country,
    this.phone,
    this.isDefault = false,
  });

  String get fullAddress => '$street, $city, $state $zip, $country';

  factory Address.fromJson(Map<String, dynamic> json) {
    return Address(
      id: json['id'] as String,
      label: json['label'] as String? ?? 'Home',
      name: json['name'] as String,
      street: json['street'] as String,
      city: json['city'] as String,
      state: json['state'] as String,
      zip: json['zip'] as String,
      country: json['country'] as String? ?? 'US',
      phone: json['phone'] as String?,
      isDefault: json['isDefault'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'label': label,
      'name': name,
      'street': street,
      'city': city,
      'state': state,
      'zip': zip,
      'country': country,
      if (phone != null) 'phone': phone,
      'isDefault': isDefault,
    };
  }
}

class ProfileRepository {
  final ApiClient _apiClient;

  ProfileRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<UserProfile> getProfile() async {
    final response = await _apiClient.get('/profile');
    return UserProfile.fromJson(response.data as Map<String, dynamic>);
  }

  Future<UserProfile> updateProfile({
    String? firstName,
    String? lastName,
    String? phone,
    String? avatarUrl,
  }) async {
    final response = await _apiClient.put('/profile', data: {
      if (firstName != null) 'firstName': firstName,
      if (lastName != null) 'lastName': lastName,
      if (phone != null) 'phone': phone,
      if (avatarUrl != null) 'avatarUrl': avatarUrl,
    });
    return UserProfile.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<Address>> getAddresses() async {
    final response = await _apiClient.get('/profile/addresses');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => Address.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Address> addAddress(Address address) async {
    final response = await _apiClient.post(
      '/profile/addresses',
      data: address.toJson(),
    );
    return Address.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Address> updateAddress(String id, Address address) async {
    final response = await _apiClient.put(
      '/profile/addresses/$id',
      data: address.toJson(),
    );
    return Address.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> deleteAddress(String id) async {
    await _apiClient.delete('/profile/addresses/$id');
  }
}
