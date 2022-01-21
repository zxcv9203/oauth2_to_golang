# oauth2란?
구글, 페이스북, 카카오 등에서 제공하는 인증 서버를 통해
회원 정보를 인증하고 Access Token을 발급받기 위한 표준 프로토콜입니다.

발급받은 Access Token을 이용하여 인증 받은 곳(구글, 페이스북, 카카오 등)
의 API 서비스를 이용할 수 있게 됩니다.

## oauth2 용어

Access Token : Authorization Server로 부터 발급 받은 인증 토큰 으로 Resource Server에 전달하여 서비스를 제공받을 수 있습니다.

Refresh Token : Access Token이 만료된 경우 클라이언트가 Refresh Token을 이용하여 새로운 Access Token으로 교환하는데 사용됩니다.

Resource owner : Resource server로 부터 계정을 소유하고 있는 사용자를 의미합니다.

Client : 구글, 페이스북, 카카오 등의 API 서비스를 이용하는 제 3의 서비스를 의미합니다.

Authorization Server(권한 서버) : 권한을 관리해주는 서버, Access Token, Refresh Token을 발급, 재발급 해주는 역할을 합니다.

Resource Server : OAuth 서비스를 제공하고, 자원을 관리하는 서버입니다.

## OAuth2의 인증방식

1. Authorization Code Grant (권한 부여 승인 코드 방식)

	권한 부여 승인을 위해 자체 생성한 Authorization Code를 전달하는 방식으로 기본이 되는 방식입니다.
	간편 로그인 기능에서 사용되는 방식으로 클라이언트가 사용자를 대신하여 특정 자원에 접근을 요청할 때 사용되는 방식입니다.
	보통 타사의 클라이언트에게 보호된 자원을 제공하기 위한 인증에 사용됩니다.
```
     +----------+
     | Resource |
     |   Owner  |
     |          |
     +----------+
          ^
          |
         (B)
     +----|-----+          Client Identifier      +---------------+
     |         -+----(A)-- & Redirection URI ---->|               |
     |  User-   |                                 | Authorization |
     |  Agent  -+----(B)-- User authenticates --->|     Server    |
     |          |                                 |               |
     |         -+<---(C)-- Authorization Code ----|               |
     +-|----|---+                                 +---------------+
       |    |                                         ^      v
      (A)  (C)                                        |      |
       |    |                                         |      |
       ^    v                                         |      |
     +---------+                                      |      |
     |         |>---(D)-- Authorization Code ---------'      |
     |  Client |          & Redirection URI                  |
     |         |                                             |
     |         |<---(E)----- Access Token -------------------'
     +---------+       (w/ Optional Refresh Token)
```
2. Implicit Grant (암묵적 승인 방식)

	Authorization Code Grant에서 인증 코드 교환과정을 제외하고 바로 access token을 발급받는 방법입니다.

	특별히 안전한 저장공간이 없는 JavaScript SPA(Single Page Application)에 사용하기 위해 만들어졌지만 권장되는 방식이 아닙니다.
```
     +----------+
     | Resource |
     |  Owner   |
     |          |
     +----------+
          ^
          |
         (B)
     +----|-----+          Client Identifier     +---------------+
     |         -+----(A)-- & Redirection URI --->|               |
     |  User-   |                                | Authorization |
     |  Agent  -|----(B)-- User authenticates -->|     Server    |
     |          |                                |               |
     |          |<---(C)--- Redirection URI ----<|               |
     |          |          with Access Token     +---------------+
     |          |            in Fragment
     |          |                                +---------------+
     |          |----(D)--- Redirection URI ---->|   Web-Hosted  |
     |          |          without Fragment      |     Client    |
     |          |                                |    Resource   |
     |     (F)  |<---(E)------- Script ---------<|               |
     |          |                                +---------------+
     +-|--------+
       |    |
      (A)  (G) Access Token
       |    |
       ^    v
     +---------+
     |         |
     |  Client |
     |         |
     +---------+
```

3. Resource Owner Password Credentials Grant (리소스 소유자 암호 자격 증명 방식)

	외부의 프로그램이나 타 사의 서비스를 접근하는 것이 아닌 본인의 서비스의 애플리케이션에서 사용하는 인증방법입니다.

	사용자 이름과 비밀번호를 사용하여 access token을 얻습니다.
```
     +----------+
     | Resource |
     |  Owner   |
     |          |
     +----------+
          v
          |    Resource Owner
         (A) Password Credentials
          |
          v
     +---------+                                  +---------------+
     |         |>--(B)---- Resource Owner ------->|               |
     |         |         Password Credentials     | Authorization |
     | Client  |                                  |     Server    |
     |         |<--(C)---- Access Token ---------<|               |
     |         |    (w/ Optional Refresh Token)   |               |
     +---------+                                  +---------------+
```

4. Client Credentials Grant (클라이언트 자격증명 승인 방식)

	단순히 클라이언트가 인증서버로 요청을 보내면 액세스 토큰을 반환하는 방식입니다.
```
     +---------+                                  +---------------+
     |         |                                  |               |
     |         |>--(A)- Client Authentication --->| Authorization |
     | Client  |                                  |     Server    |
     |         |<--(B)---- Access Token ---------<|               |
     |         |                                  |               |
     +---------+                                  +---------------+
```

## OAuth2의 통신 흐름
```

     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+
```

