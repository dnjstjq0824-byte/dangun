class User {
  final String id;
  final String email;
  final String nickname;
  final double mannerTemperature;

  User({
    required this.id,
    required this.email,
    required this.nickname,
    this.mannerTemperature = 36.5,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      nickname: json['nickname'] ?? '',
      mannerTemperature: (json['profile']?['manner_temperature'] as num?)?.toDouble() ?? 36.5,
    );
  }
}

class Listing {
  final String id;
  final String sellerUserId;
  final String sellerNickname;
  final double sellerMannerTemp;
  final String neighborhoodName;
  final String categoryName;
  final String title;
  final String description;
  final String condition;
  final int priceMinor;
  final String currencyCode;
  final bool isFree;
  final bool isNegotiable;
  final String listingStatus;
  final List<String> images;
  final int favoriteCount;
  final int chatCount;
  final bool isFavorited;

  Listing({
    required this.id,
    required this.sellerUserId,
    required this.sellerNickname,
    this.sellerMannerTemp = 36.5,
    required this.neighborhoodName,
    required this.categoryName,
    required this.title,
    required this.description,
    required this.condition,
    required this.priceMinor,
    required this.currencyCode,
    required this.isFree,
    required this.isNegotiable,
    required this.listingStatus,
    required this.images,
    this.favoriteCount = 0,
    this.chatCount = 0,
    this.isFavorited = false,
  });

  factory Listing.fromJson(Map<String, dynamic> json) {
    return Listing(
      id: json['id'] ?? '',
      sellerUserId: json['seller_user_id'] ?? '',
      sellerNickname: json['seller_nickname'] ?? '판매자',
      sellerMannerTemp: (json['seller_manner_temp'] as num?)?.toDouble() ?? 36.5,
      neighborhoodName: json['neighborhood_name'] ?? '청라동',
      categoryName: json['category_name'] ?? '디지털기기',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      condition: json['condition'] ?? 'GOOD',
      priceMinor: json['price_minor'] ?? 0,
      currencyCode: json['currency_code'] ?? 'KRW',
      isFree: json['is_free'] ?? false,
      isNegotiable: json['is_negotiable'] ?? true,
      listingStatus: json['listing_status'] ?? 'ACTIVE',
      images: (json['images'] as List<dynamic>?)?.map((e) => e.toString()).toList() ?? [],
      favoriteCount: json['favorite_count'] ?? 0,
      chatCount: json['chat_count'] ?? 0,
      isFavorited: json['is_favorited'] ?? false,
    );
  }
}

class Conversation {
  final String id;
  final String listingId;
  final String listingTitle;
  final int listingPrice;
  final String listingImage;
  final String listingStatus;
  final String buyerUserId;
  final String buyerNickname;
  final String sellerUserId;
  final String sellerNickname;
  final String lastMessage;
  final int unreadCount;

  Conversation({
    required this.id,
    required this.listingId,
    required this.listingTitle,
    required this.listingPrice,
    required this.listingImage,
    required this.listingStatus,
    required this.buyerUserId,
    required this.buyerNickname,
    required this.sellerUserId,
    required this.sellerNickname,
    required this.lastMessage,
    this.unreadCount = 0,
  });

  factory Conversation.fromJson(Map<String, dynamic> json) {
    return Conversation(
      id: json['id'] ?? '',
      listingId: json['listing_id'] ?? '',
      listingTitle: json['listing_title'] ?? '',
      listingPrice: json['listing_price'] ?? 0,
      listingImage: json['listing_image'] ?? '',
      listingStatus: json['listing_status'] ?? 'ACTIVE',
      buyerUserId: json['buyer_user_id'] ?? '',
      buyerNickname: json['buyer_nickname'] ?? '구매자',
      sellerUserId: json['seller_user_id'] ?? '',
      sellerNickname: json['seller_nickname'] ?? '판매자',
      lastMessage: json['last_message'] ?? '',
      unreadCount: json['unread_count'] ?? 0,
    );
  }
}

class ChatMessage {
  final String id;
  final String conversationId;
  final String senderUserId;
  final int sequenceNo;
  final String messageType;
  final String body;
  final DateTime createdAt;

  ChatMessage({
    required this.id,
    required this.conversationId,
    required this.senderUserId,
    required this.sequenceNo,
    required this.messageType,
    required this.body,
    required this.createdAt,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      id: json['id'] ?? '',
      conversationId: json['conversation_id'] ?? '',
      senderUserId: json['sender_user_id'] ?? '',
      sequenceNo: json['sequence_no'] ?? 0,
      messageType: json['message_type'] ?? 'TEXT',
      body: json['body'] ?? '',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class CommunityPost {
  final String id;
  final String authorUserId;
  final String authorNickname;
  final String neighborhoodName;
  final String category;
  final String title;
  final String body;
  final int likeCount;
  final int commentCount;
  final bool isLiked;

  CommunityPost({
    required this.id,
    required this.authorUserId,
    required this.authorNickname,
    required this.neighborhoodName,
    required this.category,
    required this.title,
    required this.body,
    this.likeCount = 0,
    this.commentCount = 0,
    this.isLiked = false,
  });

  factory CommunityPost.fromJson(Map<String, dynamic> json) {
    return CommunityPost(
      id: json['id'] ?? '',
      authorUserId: json['author_user_id'] ?? '',
      authorNickname: json['author_nickname'] ?? '이웃',
      neighborhoodName: json['neighborhood_name'] ?? '청라동',
      category: json['category'] ?? '일상',
      title: json['title'] ?? '',
      body: json['body'] ?? '',
      likeCount: json['like_count'] ?? 0,
      commentCount: json['comment_count'] ?? 0,
      isLiked: json['is_liked'] ?? false,
    );
  }
}

