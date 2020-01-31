# Rest API Secure Auth

# About
- golangでセキュアな認証システムを実装した
- 認証には goth を使っており、今回はTwitterログインのみを実装した
- SPAと組み合わせて使えるような REST API を意識して開発した
- 自動でログインされるようにcookieを保持している
- DB は Mock で作っており、メモリー上に情報を保持している
- 確認のためにここで確認できるようにカバーするべき項目を列挙している

# Secure API
- [x] 基本的に`gorilla`関連のpackageは様々な機能を小分けで提供してくれているため必要な物を組み合わせて安全に開発を行える

### セッション・Cookie
- [x] セッション情報は`Cookie`または`JWT`に保存する
  - **Cookie** 場合にもよるがこっちの方が良さそう?
    - `Cookie`を使う場合は`httpOnly`や`secure`などのオプションから安全な設定を追加できるため、`XSS`からの攻撃を防ぐことできるが、`CSRF`は自前で実装して防ぐ必要がある
    - 安全に実装するなら`Cookie`のような気がする
  - **JWT**
    - `JWT`はブラウザの`LocalStrage`に保存することができ、`JWT`に情報を持たせることができるので`Server`をステートレスに保つことができる
    - `LocalStrage`は`Same-Origin`の場合のみでしか、I/O処理を行うことができないため`CSRF`の問題はないが、`XSS`によって情報を抜き取られる可能性がある
    - `XSS`を100%含まないと言い切れるサイトはない?(https://techracho.bpsinc.jp/hachi8833/2019_10_09/80851)
- [x] Cookie の扱いに気をつける
  - Cookie を信用しすぎない設計にする
  - ユーザー情報の編集などの個人情報の編集には必ず Password を求めるようにする
- [x] 予測可能な情報をSessionIDに指定してはいけない
  - IDや日時など推測可能な値をSessionIDに指定すると、簡単に推測されて不正に情報をとられ、ログインなどの処理が可能になる
  - また推測可能な値をハッシュ化して指定することもよくない(時間をかければ推測される可能性があるため)
  - そのためSessionIDの生成には、言語が指定指定しているプログラムやフレームワークなどの機能を使うと良い
  - Goでは`gorilla/sessions`を使うと楽
 - [x] SessionIDの固定化攻撃を防ぐ
  - SessionIDの固定化攻撃とは攻撃者が被害者に対して、SessionIDをなんらの方法で指定することにより、指定されたSessionIDで被害者がログインすると、攻撃者は指定したSessionIDにより、ログイン状態となる脆弱性である。
  - また、ログインしていなくても入力した情報をSessionに逐一保存している場合、固定化攻撃によりSessionIDを指定されると、そのSessionに情報が蓄積されることで攻撃者に情報が抜き取られる可能性がある
  - 多くの場合は心配ないが、セッションアダプションというセッションを外部から指定できるような機能を持っている言語(PHPなど)で起きやすい。しかし、基本的にこれらの機能はデフォルトでfalseになっているはずなので心配はいらないはずである。
  - Cookie に SessionID を保存する(URLに保存しない)
  - 認証成功後に SessionID を変更する(変更できない場合はTokenによりSessionIDの認証を行う)
  - 認証前に機密情報をSessionに保存しない
- [x] Cookieの`httpOnly`と`secure`を`true`にする(`httpOnly`はJSからアクセス不可能にするためで、`secure`は`https`でのみCookieを扱うことを指定する)
  - 開発の段階で`secure`を`true`にしていると`localhost`で使用できない可能性があるため、開発時は`false`で良い(公開する時には`true`にすること)

### CORS
- [x] `CORS`をちゃんと設定する([オリジン間リソース共有 (CORS)](https://developer.mozilla.org/ja/docs/Web/HTTP/CORS))
  - `Access-Control-Allow-Origin` ... 許可するOriginを指定する(デフォルトは同じOriginが指定される)
  - `Access-Control-Allow-Methods` ... 許可する`HTTP Method`を指定する(`GET, POST, OPTIONS, HEAD`など)
  - `Access-Control-Allow-Headers` ... 許可するヘッダーを指定する。プリフライトリクエストのレスポンスで使用される。(`Content-Type, Authorization`など)
  - `Access-Control-Allow-Credentials` ... 資格情報が必要なリクエストに対して、レスポンスを開示するかどうか(普通は何も指定しなくて良い) 
  - `Access-Control-Max-Age` ... プリフライトリクエストを何度も呼ぶのはオーバーヘッドになるので、このヘッダーに時間を指定することでキャッシュさせることができる
  - 上記の`CORS`をしっかり設定した上で`CSRF Token`をレスポンスする

### Preflight Request
- `Preflight Request`とは、主にJSからのPOSTなどの副作用を保つメソッドに対するリクエストを行う時に、安全なリクエストを送るために事前にリクエストされる通信である
- `Preflight Request`によって`Access-Control-*`のHeader情報が検証されることで安全なリクエストを行うことができる

### CSRF
- [x] `CSRF`対策をする
  - [gorilla/csrf](https://github.com/gorilla/csrf#javascript-applications)を使うと楽
  - CSRFの必要性([これで完璧！今さら振り返る CSRF 対策と同一オリジンポリシーの基礎](https://qiita.com/mpyw/items/0595f07736cfa5b1f50c), [gorilla/csrf で安全なWebフォームを作る](http://matope.hatenablog.com/entry/2019/06/05/144435))
  - cookieに CSRF Token を保存しておき、Client に Response する
  - Client では受け取った Token を Request に含めて送信する
  - `GET, OPTIONS, HEAD, TRACE`はCSRFの検証をする必要がないはず(データの変更を行うような処理を含まないため)
  - `JWT`を使うことでステートレスなCSRF対策ができる(https://qiita.com/kaiinui/items/21ec7cc8a1130a1a103a)
  - CSRF対策として`Preflight Request`を使う方法もあるが、`CSRF Token`を発行していれば、Same Origin であることの検証は可能なので必要ない

### X-Frame-Options
- クリックジャッキングなどの脆弱性対策として必要
- クリックジャッキングは、攻撃者が作成した偽サイトに`iframe`を使ってTwitterなどの一般的なサイトを表示し、その一般的なサイトの投稿ボタンなどの上に罠サイトへのリンクや、攻撃用のプログラムを含めておき、実行させるというものである。
- この攻撃を防ぐためには、Originまたは指定されたOrigin以外のサイトでは`iframe`を使用できないようにする必要がある。
- それが`X-Frame-Options`である
- これを含めることで、知らないサイトから`iframe`を使って自分のサイトを表示される事はなくなるため、脆弱性を防ぐことができる
- 実際にTwitterを`iframe`から参照しようとすると`Refused to display 'https://twitter.com/' in a frame because it set 'X-Frame-Options' to 'deny'`というエラーが出力される
- 参考: https://qiita.com/mejileben/items/39d897757d5c3a904721

### X-XSS-Protection
- XSSを検知した時に無害な出力に変換する
- 各ブラウザでデフォルトでは有効になっているが、ユーザーの設定で無効にできるため、`X-XSS-Protection`を設定することで矯正させることができる

### X-Content-Type-Options
- MIME Type が変更されないように強制するために付与する
- https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/X-Content-Type-Options

### その他
- [x] SQL Injection
- [x] Passwordなどの見られてはいけない重要な情報を暗号化してからDBに保存する
