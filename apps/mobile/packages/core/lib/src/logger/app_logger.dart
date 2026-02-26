import 'package:logger/logger.dart';

/// Application-wide logger wrapper built on top of the `logger` package.
///
/// Provides convenient logging methods with consistent formatting and
/// optional tag support for filtering log output by feature or module.
class AppLogger {
  /// Internal logger instance with a pretty printer for development.
  final Logger _logger;

  /// Creates an [AppLogger] with an optional custom [Logger] instance.
  ///
  /// If no logger is provided, a default one is created with [PrettyPrinter].
  AppLogger({Logger? logger})
      : _logger = logger ??
            Logger(
              printer: PrettyPrinter(
                methodCount: 0,
                errorMethodCount: 5,
                lineLength: 80,
                colors: true,
                printEmojis: true,
                dateTimeFormat: DateTimeFormat.onlyTimeAndSinceStart,
              ),
            );

  /// Logs a verbose / trace-level message.
  void verbose(dynamic message, {String? tag}) {
    _logger.t(_formatMessage(message, tag));
  }

  /// Logs a debug-level message.
  void debug(dynamic message, {String? tag}) {
    _logger.d(_formatMessage(message, tag));
  }

  /// Logs an info-level message.
  void info(dynamic message, {String? tag}) {
    _logger.i(_formatMessage(message, tag));
  }

  /// Logs a warning-level message.
  void warning(dynamic message, {String? tag}) {
    _logger.w(_formatMessage(message, tag));
  }

  /// Logs an error-level message with optional [error] and [stackTrace].
  void error(
    dynamic message, {
    String? tag,
    dynamic error,
    StackTrace? stackTrace,
  }) {
    _logger.e(
      _formatMessage(message, tag),
      error: error,
      stackTrace: stackTrace,
    );
  }

  /// Logs a fatal / what-a-terrible-failure message.
  void fatal(
    dynamic message, {
    String? tag,
    dynamic error,
    StackTrace? stackTrace,
  }) {
    _logger.f(
      _formatMessage(message, tag),
      error: error,
      stackTrace: stackTrace,
    );
  }

  /// Prefixes the message with an optional [tag] for easy filtering.
  String _formatMessage(dynamic message, String? tag) {
    if (tag != null) {
      return '[$tag] $message';
    }
    return message.toString();
  }
}
