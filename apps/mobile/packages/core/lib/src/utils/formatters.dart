import 'package:intl/intl.dart';

/// Collection of formatting utilities for prices, dates, and other
/// display values used throughout the application.
class Formatters {
  Formatters._();

  // ---------------------------------------------------------------------------
  // Price / Currency
  // ---------------------------------------------------------------------------

  /// Converts a price stored in **cents** to a formatted dollar string.
  ///
  /// Example: `formatPrice(1299)` returns `"$12.99"`.
  ///
  /// [currency] defaults to `'USD'` and controls the currency symbol via
  /// the `intl` package.
  static String formatPrice(int cents, {String currency = 'USD'}) {
    final dollars = cents / 100.0;
    final formatter = NumberFormat.currency(
      locale: 'en_US',
      symbol: _currencySymbol(currency),
      decimalDigits: 2,
    );
    return formatter.format(dollars);
  }

  /// Returns a compact price string without trailing zeros when the price
  /// is a whole dollar amount.
  ///
  /// Example: `formatPriceCompact(1200)` returns `"$12"`,
  /// while `formatPriceCompact(1299)` returns `"$12.99"`.
  static String formatPriceCompact(int cents, {String currency = 'USD'}) {
    final dollars = cents / 100.0;
    if (dollars == dollars.truncateToDouble()) {
      final formatter = NumberFormat.currency(
        locale: 'en_US',
        symbol: _currencySymbol(currency),
        decimalDigits: 0,
      );
      return formatter.format(dollars);
    }
    return formatPrice(cents, currency: currency);
  }

  /// Returns the currency symbol for a given ISO 4217 currency code.
  static String _currencySymbol(String currencyCode) {
    switch (currencyCode.toUpperCase()) {
      case 'USD':
        return '\$';
      case 'EUR':
        return '\u20AC';
      case 'GBP':
        return '\u00A3';
      case 'JPY':
        return '\u00A5';
      case 'INR':
        return '\u20B9';
      default:
        return currencyCode;
    }
  }

  // ---------------------------------------------------------------------------
  // Date / Time
  // ---------------------------------------------------------------------------

  /// Formats a [DateTime] as `"Jan 15, 2024"`.
  static String formatDate(DateTime date) {
    return DateFormat.yMMMd('en_US').format(date);
  }

  /// Formats a [DateTime] as `"January 15, 2024"`.
  static String formatDateLong(DateTime date) {
    return DateFormat.yMMMMd('en_US').format(date);
  }

  /// Formats a [DateTime] as `"01/15/2024"`.
  static String formatDateShort(DateTime date) {
    return DateFormat('MM/dd/yyyy').format(date);
  }

  /// Formats a [DateTime] as `"3:30 PM"`.
  static String formatTime(DateTime date) {
    return DateFormat.jm('en_US').format(date);
  }

  /// Formats a [DateTime] as `"Jan 15, 2024 3:30 PM"`.
  static String formatDateTime(DateTime date) {
    return DateFormat('MMM d, yyyy h:mm a', 'en_US').format(date);
  }

  /// Returns a human-readable relative time string such as
  /// `"just now"`, `"5 minutes ago"`, `"2 hours ago"`, or `"3 days ago"`.
  static String formatRelativeTime(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inSeconds < 60) {
      return 'just now';
    } else if (difference.inMinutes < 60) {
      final minutes = difference.inMinutes;
      return '$minutes ${minutes == 1 ? 'minute' : 'minutes'} ago';
    } else if (difference.inHours < 24) {
      final hours = difference.inHours;
      return '$hours ${hours == 1 ? 'hour' : 'hours'} ago';
    } else if (difference.inDays < 7) {
      final days = difference.inDays;
      return '$days ${days == 1 ? 'day' : 'days'} ago';
    } else if (difference.inDays < 30) {
      final weeks = (difference.inDays / 7).floor();
      return '$weeks ${weeks == 1 ? 'week' : 'weeks'} ago';
    } else if (difference.inDays < 365) {
      final months = (difference.inDays / 30).floor();
      return '$months ${months == 1 ? 'month' : 'months'} ago';
    } else {
      final years = (difference.inDays / 365).floor();
      return '$years ${years == 1 ? 'year' : 'years'} ago';
    }
  }

  // ---------------------------------------------------------------------------
  // Numbers
  // ---------------------------------------------------------------------------

  /// Formats a large number with commas: `1234567` becomes `"1,234,567"`.
  static String formatNumber(num value) {
    return NumberFormat('#,##0', 'en_US').format(value);
  }

  /// Formats a number in compact form: `1200` becomes `"1.2K"`.
  static String formatCompactNumber(num value) {
    return NumberFormat.compact(locale: 'en_US').format(value);
  }

  /// Formats a decimal as a percentage: `0.856` becomes `"85.6%"`.
  static String formatPercentage(double value, {int decimalDigits = 1}) {
    return '${(value * 100).toStringAsFixed(decimalDigits)}%';
  }

  // ---------------------------------------------------------------------------
  // Order
  // ---------------------------------------------------------------------------

  /// Formats an order number with a leading hash: `"ORD-00001234"`.
  static String formatOrderNumber(String orderNumber) {
    return '#$orderNumber';
  }
}
