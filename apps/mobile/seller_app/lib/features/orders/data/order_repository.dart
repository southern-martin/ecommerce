/// Represents an item within a seller order.
class OrderItem {
  final String id;
  final String productId;
  final String productName;
  final String? imageUrl;
  final int quantity;
  final double price;
  final double total;

  const OrderItem({
    required this.id,
    required this.productId,
    required this.productName,
    this.imageUrl,
    required this.quantity,
    required this.price,
    required this.total,
  });
}

/// Represents a customer's shipping address.
class CustomerInfo {
  final String name;
  final String email;
  final String phone;
  final String address;
  final String city;
  final String state;
  final String zipCode;

  const CustomerInfo({
    required this.name,
    required this.email,
    required this.phone,
    required this.address,
    required this.city,
    required this.state,
    required this.zipCode,
  });

  String get fullAddress => '$address, $city, $state $zipCode';
}

/// Represents a seller order.
class SellerOrder {
  final String id;
  final String orderNumber;
  final String customerName;
  final CustomerInfo? customerInfo;
  final List<OrderItem> items;
  final double total;
  final String status;
  final String? trackingNumber;
  final DateTime createdAt;
  final DateTime updatedAt;

  const SellerOrder({
    required this.id,
    required this.orderNumber,
    required this.customerName,
    this.customerInfo,
    required this.items,
    required this.total,
    required this.status,
    this.trackingNumber,
    required this.createdAt,
    required this.updatedAt,
  });

  int get itemCount => items.fold(0, (sum, item) => sum + item.quantity);
}

/// Paginated result wrapper for orders.
class PaginatedOrders {
  final List<SellerOrder> orders;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedOrders({
    required this.orders,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });
}

/// Order statistics summary.
class OrderStats {
  final int totalOrders;
  final int pendingOrders;
  final int processingOrders;
  final int shippedOrders;
  final int deliveredOrders;
  final int cancelledOrders;

  const OrderStats({
    required this.totalOrders,
    required this.pendingOrders,
    required this.processingOrders,
    required this.shippedOrders,
    required this.deliveredOrders,
    required this.cancelledOrders,
  });
}

/// Repository for managing seller orders.
class SellerOrderRepository {
  /// Fetches paginated list of seller orders with optional status filter.
  Future<PaginatedOrders> getOrders({
    int page = 1,
    String? status,
  }) async {
    // TODO: Replace with actual API call to /seller/orders
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    final statuses = ['pending', 'processing', 'shipped', 'delivered', 'cancelled'];
    final names = ['Alice Johnson', 'Bob Smith', 'Carol Davis', 'Dan Wilson', 'Eve Brown'];

    final allOrders = List.generate(
      20,
      (i) => SellerOrder(
        id: 'order_${i + 1}',
        orderNumber: 'ORD-${1001 + i}',
        customerName: names[i % names.length],
        items: List.generate(
          (i % 3) + 1,
          (j) => OrderItem(
            id: 'item_${i}_$j',
            productId: 'prod_${j + 1}',
            productName: 'Product ${j + 1}',
            imageUrl: 'https://picsum.photos/200?random=${i * 10 + j}',
            quantity: (j % 3) + 1,
            price: 19.99 + (j * 15.0),
            total: (19.99 + (j * 15.0)) * ((j % 3) + 1),
          ),
        ),
        total: 49.99 + (i * 25.0),
        status: statuses[i % statuses.length],
        trackingNumber: i % statuses.length == 2 ? 'TRK${100000 + i}' : null,
        createdAt: now.subtract(Duration(days: i, hours: i * 3)),
        updatedAt: now.subtract(Duration(days: i)),
      ),
    );

    final filtered = status != null
        ? allOrders.where((o) => o.status == status).toList()
        : allOrders;

    return PaginatedOrders(
      orders: filtered,
      totalCount: filtered.length,
      currentPage: page,
      totalPages: 1,
    );
  }

  /// Fetches a single order by ID.
  Future<SellerOrder> getOrderById(String id) async {
    // TODO: Replace with actual API call to /seller/orders/:id
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return SellerOrder(
      id: id,
      orderNumber: 'ORD-1001',
      customerName: 'Alice Johnson',
      customerInfo: const CustomerInfo(
        name: 'Alice Johnson',
        email: 'alice@example.com',
        phone: '+1 (555) 123-4567',
        address: '123 Main Street',
        city: 'Springfield',
        state: 'IL',
        zipCode: '62701',
      ),
      items: [
        const OrderItem(
          id: 'item_1',
          productId: 'prod_1',
          productName: 'Wireless Bluetooth Headphones',
          imageUrl: 'https://picsum.photos/200?random=1',
          quantity: 1,
          price: 79.99,
          total: 79.99,
        ),
        const OrderItem(
          id: 'item_2',
          productId: 'prod_2',
          productName: 'Phone Case - Clear',
          imageUrl: 'https://picsum.photos/200?random=2',
          quantity: 2,
          price: 14.99,
          total: 29.98,
        ),
      ],
      total: 109.97,
      status: 'pending',
      createdAt: now.subtract(const Duration(hours: 5)),
      updatedAt: now.subtract(const Duration(hours: 2)),
    );
  }

  /// Updates the status of an order.
  Future<SellerOrder> updateOrderStatus(
    String id,
    String status, {
    String? trackingNumber,
  }) async {
    // TODO: Replace with actual API call to /seller/orders/:id/status
    await Future.delayed(const Duration(milliseconds: 800));

    final order = await getOrderById(id);
    return SellerOrder(
      id: order.id,
      orderNumber: order.orderNumber,
      customerName: order.customerName,
      customerInfo: order.customerInfo,
      items: order.items,
      total: order.total,
      status: status,
      trackingNumber: trackingNumber ?? order.trackingNumber,
      createdAt: order.createdAt,
      updatedAt: DateTime.now(),
    );
  }

  /// Fetches order statistics summary.
  Future<OrderStats> getOrderStats() async {
    // TODO: Replace with actual API call to /seller/orders/stats
    await Future.delayed(const Duration(milliseconds: 500));

    return const OrderStats(
      totalOrders: 187,
      pendingOrders: 12,
      processingOrders: 8,
      shippedOrders: 25,
      deliveredOrders: 135,
      cancelledOrders: 7,
    );
  }
}
