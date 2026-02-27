import 'package:flutter/material.dart';

class OAuthButtons extends StatefulWidget {
  final VoidCallback? onSuccess;

  const OAuthButtons({super.key, this.onSuccess});

  @override
  State<OAuthButtons> createState() => _OAuthButtonsState();
}

class _OAuthButtonsState extends State<OAuthButtons> {
  bool _isLoading = false;

  Future<void> _handleOAuthLogin(String provider) async {
    setState(() => _isLoading = true);
    try {
      // Native SDK integration required for production:
      // Google: google_sign_in package → get idToken
      // Apple: sign_in_with_apple package → get identityToken
      // For now, show a message that platform setup is needed
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${provider[0].toUpperCase()}${provider.substring(1)} sign-in requires native SDK configuration'),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('OAuth failed: ${e.toString()}')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: OutlinedButton.icon(
            onPressed: _isLoading ? null : () => _handleOAuthLogin('google'),
            style: OutlinedButton.styleFrom(
              minimumSize: const Size(0, 52),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              side: BorderSide(color: Colors.grey.shade300),
            ),
            icon: const Icon(Icons.g_mobiledata, size: 24),
            label: const Text('Google'),
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: OutlinedButton.icon(
            onPressed: _isLoading ? null : () => _handleOAuthLogin('apple'),
            style: OutlinedButton.styleFrom(
              minimumSize: const Size(0, 52),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              side: BorderSide(color: Colors.grey.shade300),
            ),
            icon: const Icon(Icons.apple, size: 24),
            label: const Text('Apple'),
          ),
        ),
      ],
    );
  }
}
