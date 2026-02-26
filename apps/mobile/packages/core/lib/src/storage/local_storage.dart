import 'package:shared_preferences/shared_preferences.dart';

/// Wrapper around [SharedPreferences] for persisting non-sensitive,
/// lightweight key-value data such as user preferences, flags, and
/// cached settings.
class LocalStorage {
  SharedPreferences? _prefs;

  /// Initializes the underlying [SharedPreferences] instance.
  ///
  /// Must be called (and awaited) before any read/write operations.
  Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  SharedPreferences get _preferences {
    if (_prefs == null) {
      throw StateError(
        'LocalStorage has not been initialized. Call init() first.',
      );
    }
    return _prefs!;
  }

  // ---------------------------------------------------------------------------
  // String
  // ---------------------------------------------------------------------------

  /// Persists a [String] value for the given [key].
  Future<bool> setString(String key, String value) {
    return _preferences.setString(key, value);
  }

  /// Retrieves a stored [String], or `null` if no value exists.
  String? getString(String key) {
    return _preferences.getString(key);
  }

  // ---------------------------------------------------------------------------
  // Int
  // ---------------------------------------------------------------------------

  /// Persists an [int] value for the given [key].
  Future<bool> setInt(String key, int value) {
    return _preferences.setInt(key, value);
  }

  /// Retrieves a stored [int], or `null` if no value exists.
  int? getInt(String key) {
    return _preferences.getInt(key);
  }

  // ---------------------------------------------------------------------------
  // Double
  // ---------------------------------------------------------------------------

  /// Persists a [double] value for the given [key].
  Future<bool> setDouble(String key, double value) {
    return _preferences.setDouble(key, value);
  }

  /// Retrieves a stored [double], or `null` if no value exists.
  double? getDouble(String key) {
    return _preferences.getDouble(key);
  }

  // ---------------------------------------------------------------------------
  // Bool
  // ---------------------------------------------------------------------------

  /// Persists a [bool] value for the given [key].
  Future<bool> setBool(String key, bool value) {
    return _preferences.setBool(key, value);
  }

  /// Retrieves a stored [bool], or `null` if no value exists.
  bool? getBool(String key) {
    return _preferences.getBool(key);
  }

  // ---------------------------------------------------------------------------
  // String List
  // ---------------------------------------------------------------------------

  /// Persists a list of [String] values for the given [key].
  Future<bool> setStringList(String key, List<String> value) {
    return _preferences.setStringList(key, value);
  }

  /// Retrieves a stored list of [String], or `null` if no value exists.
  List<String>? getStringList(String key) {
    return _preferences.getStringList(key);
  }

  // ---------------------------------------------------------------------------
  // Generic helpers
  // ---------------------------------------------------------------------------

  /// Returns `true` if a value has been stored for the given [key].
  bool containsKey(String key) {
    return _preferences.containsKey(key);
  }

  /// Removes the value associated with [key].
  Future<bool> remove(String key) {
    return _preferences.remove(key);
  }

  /// Removes **all** stored values. Use with caution.
  Future<bool> clear() {
    return _preferences.clear();
  }

  /// Returns all stored keys.
  Set<String> getKeys() {
    return _preferences.getKeys();
  }
}
