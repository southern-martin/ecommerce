import 'package:flutter/material.dart';

class SearchSuggestions extends StatelessWidget {
  final List<String> suggestions;
  final String query;
  final ValueChanged<String> onSuggestionTap;

  const SearchSuggestions({
    super.key,
    required this.suggestions,
    required this.query,
    required this.onSuggestionTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: suggestions.length,
      itemBuilder: (context, index) {
        final suggestion = suggestions[index];
        return ListTile(
          leading: const Icon(Icons.search),
          title: _buildHighlightedText(context, suggestion),
          trailing: const Icon(Icons.north_west, size: 16),
          onTap: () => onSuggestionTap(suggestion),
        );
      },
    );
  }

  Widget _buildHighlightedText(BuildContext context, String suggestion) {
    final theme = Theme.of(context);
    final lowerSuggestion = suggestion.toLowerCase();
    final lowerQuery = query.toLowerCase();
    final matchIndex = lowerSuggestion.indexOf(lowerQuery);

    if (matchIndex < 0 || query.isEmpty) {
      return Text(suggestion);
    }

    final before = suggestion.substring(0, matchIndex);
    final match = suggestion.substring(matchIndex, matchIndex + query.length);
    final after = suggestion.substring(matchIndex + query.length);

    return RichText(
      text: TextSpan(
        style: theme.textTheme.bodyLarge?.copyWith(
          color: theme.textTheme.bodyLarge?.color,
        ),
        children: [
          if (before.isNotEmpty) TextSpan(text: before),
          TextSpan(
            text: match,
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
          if (after.isNotEmpty) TextSpan(text: after),
        ],
      ),
    );
  }
}
