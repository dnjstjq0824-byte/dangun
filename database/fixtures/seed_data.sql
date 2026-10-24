-- BN NEAR Seed Fixtures
-- Sample Data for Local Testing & Initial Launch Staging

-- 1. COUNTRIES
INSERT INTO countries (code, name, default_currency, default_locale, status) VALUES
('KR', 'South Korea', 'KRW', 'ko-KR', 'ACTIVE'),
('US', 'United States', 'USD', 'en-US', 'BETA'),
('JP', 'Japan', 'JPY', 'ja-JP', 'BETA')
ON CONFLICT (code) DO NOTHING;

-- 2. CITIES
INSERT INTO cities (id, country_code, name, timezone, status) VALUES
('a0000000-0000-0000-0000-000000000001', 'KR', 'Incheon', 'Asia/Seoul', 'ACTIVE'),
('a0000000-0000-0000-0000-000000000002', 'KR', 'Seoul', 'Asia/Seoul', 'ACTIVE'),
('a0000000-0000-0000-0000-000000000003', 'JP', 'Tokyo', 'Asia/Tokyo', 'BETA'),
('a0000000-0000-0000-0000-000000000004', 'US', 'New York', 'America/New_York', 'BETA')
ON CONFLICT (id) DO NOTHING;

-- 3. NEIGHBORHOODS (with valid polygon boundaries and centroids)
-- Cheongna-dong (Incheon Seo-gu)
INSERT INTO neighborhoods (id, city_id, slug, name, boundary_geom, public_center_geom, status) VALUES
(
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'cheongna-dong',
    '청라동 (Cheongna-dong)',
    ST_GeomFromText('POLYGON((126.6300 37.5200, 126.6700 37.5200, 126.6700 37.5500, 126.6300 37.5500, 126.6300 37.5200))', 4326),
    ST_GeomFromText('POINT(126.6500 37.5350)', 4326),
    'ACTIVE'
),
-- Songdo-dong (Incheon Yeonsu-gu)
(
    'b0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    'songdo-dong',
    '송도동 (Songdo-dong)',
    ST_GeomFromText('POLYGON((126.6200 37.3600, 126.6700 37.3600, 126.6700 37.4100, 126.6200 37.4100, 126.6200 37.3600))', 4326),
    ST_GeomFromText('POINT(126.6450 37.3850)', 4326),
    'ACTIVE'
),
-- Yeoksam-dong (Seoul Gangnam-gu)
(
    'b0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000002',
    'yeoksam-dong',
    '역삼동 (Yeoksam-dong)',
    ST_GeomFromText('POLYGON((127.0250 37.4950, 127.0450 37.4950, 127.0450 37.5100, 127.0250 37.5100, 127.0250 37.4950))', 4326),
    ST_GeomFromText('POINT(127.0350 37.5020)', 4326),
    'ACTIVE'
),
-- Yeonnam-dong (Seoul Mapo-gu)
(
    'b0000000-0000-0000-0000-000000000004',
    'a0000000-0000-0000-0000-000000000002',
    'yeonnam-dong',
    '연남동 (Yeonnam-dong)',
    ST_GeomFromText('POLYGON((126.9150 37.5600, 126.9300 37.5600, 126.9300 37.5720, 126.9150 37.5720, 126.9150 37.5600))', 4326),
    ST_GeomFromText('POINT(126.9220 37.5660)', 4326),
    'ACTIVE'
)
ON CONFLICT (id) DO NOTHING;

-- 4. CATEGORIES
INSERT INTO categories (id, slug, name, icon_name, sort_order) VALUES
('c0000000-0000-0000-0000-000000000001', 'digital-tech', '디지털기기 (Digital & Tech)', 'smartphone', 1),
('c0000000-0000-0000-0000-000000000002', 'home-appliances', '생활가전 (Appliances)', 'tv', 2),
('c0000000-0000-0000-0000-000000000003', 'furniture-interior', '가구/인테리어 (Furniture)', 'home', 3),
('c0000000-0000-0000-0000-000000000004', 'baby-kids', '유아동 (Baby & Kids)', 'heart', 4),
('c0000000-0000-0000-0000-000000000005', 'fashion-clothing', '의류/잡화 (Fashion)', 'shopping-bag', 5),
('c0000000-0000-0000-0000-000000000006', 'sports-leisure', '스포츠/레저 (Sports)', 'activity', 6),
('c0000000-0000-0000-0000-000000000007', 'books-tickets', '도서/음반/티켓 (Books & Media)', 'book', 7),
('c0000000-0000-0000-0000-000000000008', 'pet-supplies', '반려동물용품 (Pet Supplies)', 'smile', 8),
('c0000000-0000-0000-0000-000000000009', 'plants', '식물 (Plants)', 'feather', 9),
('c0000000-0000-0000-0000-000000000010', 'local-services', '동네서비스 (Local Services)', 'tool', 10)
ON CONFLICT (id) DO NOTHING;

-- 5. TEST USERS
-- Passwords hashed with bcrypt (cost 10): $2a$10$7EqJtq98hPqEX7fNZaFWoO.8/bBvN/3pD6yGv2LwFm9eCqYI8xM4S = "Password123!"
INSERT INTO users (id, email, password_hash, nickname, status, email_verified_at) VALUES
('u0000000-0000-0000-0000-000000000001', 'admin@bn-near.com', '$2a$10$7EqJtq98hPqEX7fNZaFWoO.8/bBvN/3pD6yGv2LwFm9eCqYI8xM4S', 'NearAdmin', 'ACTIVE', NOW()),
('u0000000-0000-0000-0000-000000000002', 'seller@bn-near.com', '$2a$10$7EqJtq98hPqEX7fNZaFWoO.8/bBvN/3pD6yGv2LwFm9eCqYI8xM4S', 'David_Cheongna', 'ACTIVE', NOW()),
('u0000000-0000-0000-0000-000000000003', 'buyer@bn-near.com', '$2a$10$7EqJtq98hPqEX7fNZaFWoO.8/bBvN/3pD6yGv2LwFm9eCqYI8xM4S', 'Minji_Cheongna', 'ACTIVE', NOW()),
('u0000000-0000-0000-0000-000000000004', 'merchant@bn-near.com', '$2a$10$7EqJtq98hPqEX7fNZaFWoO.8/bBvN/3pD6yGv2LwFm9eCqYI8xM4S', 'MasterKim_Bakery', 'ACTIVE', NOW())
ON CONFLICT (id) DO NOTHING;

-- Profiles
INSERT INTO user_profiles (user_id, bio, public_display_name, manner_temperature) VALUES
('u0000000-0000-0000-0000-000000000001', 'BN NEAR Official Super Admin', 'Platform Guardian', 50.0),
('u0000000-0000-0000-0000-000000000002', 'Tech enthusiast in Cheongna. Fast responses!', 'David K.', 38.2),
('u0000000-0000-0000-0000-000000000003', 'Resident of Cheongna Lake Park neighborhood', 'Minji P.', 37.0),
('u0000000-0000-0000-0000-000000000004', 'Organic Sourdough & Artisan Coffee', 'Chef Kim', 42.5)
ON CONFLICT (user_id) DO NOTHING;

-- Verified User Neighborhoods
INSERT INTO user_neighborhoods (user_id, neighborhood_id, is_primary, verification_status, verified_at, expires_at) VALUES
('u0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', TRUE, 'VERIFIED', NOW(), NOW() + INTERVAL '30 days'),
('u0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001', TRUE, 'VERIFIED', NOW(), NOW() + INTERVAL '30 days'),
('u0000000-0000-0000-0000-000000000004', 'b0000000-0000-0000-0000-000000000001', TRUE, 'VERIFIED', NOW(), NOW() + INTERVAL '30 days')
ON CONFLICT (user_id, neighborhood_id) DO NOTHING;

-- 6. SAMPLE LISTINGS
INSERT INTO listings (id, seller_user_id, neighborhood_id, category_id, title, description, condition, price_minor, currency_code, is_free, is_negotiable, listing_status, public_location_geom) VALUES
(
    'l0000000-0000-0000-0000-000000000001',
    'u0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    '애플 M2 맥북 에어 16GB 미드나이트 (Apple M2 MacBook Air)',
    '실사용 6개월 미만 배터리 효율 98% 상태 극상입니다. 박스 및 정품 충전기 풀세트 포함. 청라호수공원 정문 직거래 선호합니다.',
    'LIKE_NEW',
    1150000,
    'KRW',
    FALSE,
    TRUE,
    'ACTIVE',
    ST_GeomFromText('POINT(126.6520 37.5360)', 4326)
),
(
    'l0000000-0000-0000-0000-000000000002',
    'u0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000003',
    '이케아 원목 식탁 4인용 세트 (의자 4개 포함)',
    '이사로 인해 내놓습니다. 원목 오일케어 잘 되어 있으며 생활기스 약간 있습니다.',
    'GOOD',
    80000,
    'KRW',
    FALSE,
    FALSE,
    'ACTIVE',
    ST_GeomFromText('POINT(126.6480 37.5340)', 4326)
),
(
    'l0000000-0000-0000-0000-000000000003',
    'u0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000007',
    '동화책 전집 30권 무료 나눔 (아이들이 잘 봤어요)',
    '필요하신 이웃분께 나눔합니다. 청라 2동 주민센터 앞에서 전달 가능합니다.',
    'FAIR',
    0,
    'KRW',
    TRUE,
    FALSE,
    'ACTIVE',
    ST_GeomFromText('POINT(126.6550 37.5380)', 4326)
)
ON CONFLICT (id) DO NOTHING;

-- 7. LOCAL BUSINESS SAMPLE
INSERT INTO businesses (id, owner_user_id, name, category_id, country_code, city_id, neighborhood_id, description, public_address, phone_number, verification_status, status) VALUES
(
    'd0000000-0000-0000-0000-000000000001',
    'u0000000-0000-0000-0000-000000000004',
    '청라 아티장 베이커리 (Cheongna Artisan Bakery)',
    'c0000000-0000-0000-0000-000000000010',
    'KR',
    'a0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    '매일 아침 유기농 통밀로 굽는 천연발효 사워도우와 스페셜티 커피 전문점.',
    '인천 서구 청라커낼로 260번길 15 1층',
    '032-567-8901',
    'VERIFIED',
    'ACTIVE'
)
ON CONFLICT (id) DO NOTHING;

-- 8. COMMUNITY POST
INSERT INTO community_posts (id, author_user_id, neighborhood_id, category, title, body, like_count, comment_count) VALUES
(
    'p0000000-0000-0000-0000-000000000001',
    'u0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000001',
    '동네질문',
    '청라 호수공원 근처 주말에 아이 데리고 가기 좋은 소아과 추천 부탁드려요!',
    '최근에 청라로 이사 왔는데 주말 진료 잘 봐주시는 소아청소년과가 어디가 좋은지 이웃분들 조언 구합니다 :)',
    5,
    2
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO community_comments (id, post_id, author_user_id, body) VALUES
('m0000000-0000-0000-0000-000000000001', 'p0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002', '커낼웨이 수변상가 쪽 연세소아과 선생님 친절하시고 일요일 오전에도 진료해요!'),
('m0000000-0000-0000-0000-000000000002', 'p0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000004', '저희 아이도 연세소아과 다닙니다. 주말엔 대기 조금 있으니 일찍 가시는 걸 추천드려요~')
ON CONFLICT (id) DO NOTHING;

