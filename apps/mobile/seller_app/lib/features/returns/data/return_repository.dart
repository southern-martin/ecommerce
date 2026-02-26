/// Represents an item being returned.
class ReturnItem {
  final String id;
  final String productId;
  final String productName;
  final String? imageUrl;
  final int quantity;
  final double price;

  const ReturnItem({
    required this.id,
    required this.productId,
    required this.productName,
    this.imageUrl,
    required this.quantity,
    required this.price,
  });
}

/// Represents a return request from a customer.
class SellerReturn {
  final String id;
  final String returnNumber;
  final String orderId;
  final String orderNumber;
  final String customerName;
  final String reason;
  final String status; // pending, approved, rejected
  final List<ReturnItem> items;
  final String? note;
  final DateTime createdAt;
  final DateTime updatedAt;

  const SellerReturn({
    required this.id,
    required this.returnNumber,
    required this.orderId,
    required this.orderNumber,
    required this.customerName,
    required this.reason,
    required this.status,
    required this.items,
    this.note,
    required this.createdAt,
    required this.updatedAt,
  });
}

/// Paginated result wrapper for returns.
class PaginatedReturns {
  final List<SellerReturn> returns;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedReturns({
    required this.returns,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });
}

/// Repository for managing seller returns.
class SellerReturnRepository {
  /// Fetches paginated list of returns with optional status filter.
  Future<PaginatedReturns> getReturns({
    int page = 1,
    String? status,
  }) async {
    // TODO: Replace with actual API call to /seller/returns
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    final statuses = ['pending', 'approved', 'rejected'];
    final reasons = [
      'Item damaged during shipping',
      'Wrong item received',
      'Item does not match description',
      'Changed my mind',
      'Defective product',
    ];
    final names = ['Alice Johnson', 'Bob Smith', 'Carol Davis', 'Dan Wilson', 'Eve Brown'];

    final allReturns = List.generate(
      10,
      (i) => SellerReturn(
        id: 'ret_${i + 1}',
        returnNumber: 'RET-${2001 + i}',
        orderId: 'order_${i + 1}',
        orderNumber: 'ORD-${1001 + i}',
        customerName: names[i % names.length],
        reason: reasons[i % reasons.length],
        status: statuses[i % statuses.length],
        items: [
          ReturnItem(
            id: 'ret_item_${i}_1',
            productId: 'prod_${i + 1}',
            productName: 'Product ${i + 1}',
            imageUrl: 'https://picsum.photos/200?random=${i + 100}',
            quantity: 1,
            price: 29.99 + (i * 10.0),
          ),
        ],
        note: i % 3 == 1 ? 'Customer provided photos of damage' : null,
        createdAt: now.subtract(Duration(days: i * 2)),
        updatedAt: now.subtract(Duration(days: i)),
      ),
    );

    final filtered = status != null
        ? allReturns.where((r) => r.status == status).toList()
        : allReturns;

    return PaginatedReturns(
      returns: filtered,
      totalCount: filtered.length,
      currentPage: page,
      totalPages: 1,
    );
  }

  /// Fetches a single return by ID.
  Future<SellerReturn> getReturnById(String id) async {
    // TODO: Replace with actual API call to /seller/returns/:id
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return SellerReturn(
      id: id,
      returnNumber: 'RET-2001',
      orderId: 'order_1',
      orderNumber: 'ORD-1001',
      customerName: 'Alice Johnson',
      reason: 'Item damaged during shipping. The packaging was torn and the product has visible scratches on the surface.',
      status: 'pending',
      items: [
        const ReturnItem(
          id: 'ret_item_1',
          productId: 'prod_1',
          productName: 'Wireless Bluetooth Headphones',
          imageUrl: 'https://picsum.photos/200?random=101',
          quantity: 1,
          price: 79.99,
        ),
        const ReturnItem(
          id: 'ret_item_2',
          productId: 'prod_3',
          productName: 'USB-C Charging Cable',
          imageUrl: 'https://picsum.photos/200?random=103',
          quantity: 2,
          price: 12.99,
        ),
      ],
      createdAt: now.subtract(const Duration(days: 2)),
      updatedAt: now.subtract(const Duration(hours: 6)),
    );
  }

  /// Approves a return request with an optional note.
  Future<SellerReturn> approveReturn(String id, {String? note}) async {
    // TODO: Replace with actual API call to /seller/returns/:id/approve
    await Future.delayed(const Duration(milliseconds: 800));

    final original = await getReturnById(id);
    return SellerReturn(
      id: original.id,
      returnNumber: original.returnNumber,
      orderId: original.orderId,
      orderNumber: original.orderNumber,
      customerName: original.customerName,
      reason: original.reason,
      status: 'approved',
      items: original.items,
      note: note,
      createdAt: original.createdAt,
      updatedAt: DateTime.now(),
    );
  }

  /// Rejects a return request with a required note.
  Future<SellerReturn> rejectReturn(String id, String note) async {
    // TODO: Replace with actual API call to /seller/returns/:id/reject
    await Future.delayed(const Duration(milliseconds: 800));

    final original = await getReturnById(id);
    return SellerReturn(
      id: original.id,
      returnNumber: original.returnNumber,
      orderId: original.orderId,
      orderNumber: original.orderNumber,
      customerName: original.customerName,
      reason: original.reason,
      status: 'rejected',
      items: original.items,
      note: note,
      createdAt: original.createdAt,
      updatedAt: DateTime.now(),
    );
  }
}
