# レイヤードアーキテクチャー説明のサンプルプログラム

ここでは実際にサンプルプログラムを通じて具体的な実装方法を解説する。

## サンプルプログラムの使い方

### 実行環境の準備

```
$ make build
```
make build で関連コンテナがビルドされ実行可能な状態になる。  
ローカルにDocker環境 (lima, Rancher Desktop等のDocker互換環境でもOK)が必要。  

### プログラムの実行
以下の様にcurlなどで直接APIを叩けば動作の確認ができる。  
```bash
$ curl -X POST http://127.0.0.1:9999/customer -d '{"name":"Tetsuya Hori","email":"hori@fedecade.com"}' -i
HTTP/1.1 201 Created
Date: Mon, 17 Apr 2023 01:47:02 GMT
Content-Length: 0
```
登録に成功するとlayeredarch-dbコンテナのlayeredarchデータベースのcustomersテーブルにレコードが追加される。  
```bash
$ mysql -utest -ptest -h127.0.0.1 -P33306 layeredarch -e 'SELECT * FROM customers'
mysql: [Warning] Using a password on the command line interface can be insecure.
+-----+--------------+-------------------+
| sid | name         | email             |
+-----+--------------+-------------------+
|   1 | Tetsuya Hori | hori@fedecade.com |
+-----+--------------+-------------------+
```

## サンプルプログラムのユースケース

以下のユースケース(ここでは分析の文脈における意味で使っている)を想定する。

- アクターはREST APIによりこのソフトウエアを操作する
- アクターは顧客を登録する
- システムは外部のWEBサービスを使用して顧客の信用チェックを行う
- 登録される顧客の情報は、氏名, emailアドレスとする
- システムは顧客をemailアドレスにより一意に識別する

## サンプルプログラムのAPI仕様

- APIの結果によりリソースが新しく作られると言う性質からPOSTメソッドを採用する
- 既に同じemailアドレスが登録されている場合はHTTPステータスコード409を返す
- 顧客の属性に許可されない値が渡された場合はHTTPステータスコード400を返す
- 顧客が信用チェックをパスしなかった場合はHTTPステータスコード400を返す
- 登録が成功した場合はHTTPステータスコード201を返す

## サンプルプログラムで使用する言語

- Go言語

## サンプルプログラムで使用する外部ライブラリ

| カテゴリ | ライブラリ |
| --- | --- |
| AWF | https://github.com/gorilla/mux |
| HTTP Client | https://github.com/dghubble/sling |
| DB | https://github.com/jmoiron/sqlx |
| DI CONTAINER | https://github.com/sarulabs/di |
| TEST | https://github.com/stretchr/testify |

📔 **Note**  
開発スピードを優先してとりあえず堀が慣れているライブラリを使用する。  

## コードの概要

このサンプルプログラムは概ね以下の構成となっている。  

| レイヤー | Clean Architectureのレイヤー | 対象 | 役割 |
| --- | --- | --- | --- |
| Domain | Enterprise Business Rules | domain/customer/Customer | 顧客を表現するValueObject |
|  |  | domain/customer/Repository | 顧客の永続化を行うRepository |
|  |  | domain/creditscreening/CreditScreening | 顧客審査を行うDomainService |
| Usecase | Application Business Rules | usecase/customermanager/CustomerManager | 顧客操作を行うUsecase |
| Interface | Interface Adapters | repository/rdbcustomerrepository | domain/customer/Repositoryのデータベースを永続化先とする実装 |
|  |  | gateway/httpcreditscreening | domain/creditscreening/CreditScreeningのREST-APIサーバーを接続先とする実装 |
|  |  | requesthandler/RequestHandler | HTTPリクエストに応答するハンドラー |
|  |  | requesthandler/customerpost/ | Method: POST,  Path: /customerのリククエストに応答するハンドラーの実装 |
| Infrastructure | Frameworks & Drivers | waf | WEBアプリケーションフレームワークと関連コード |
|  |  | httpclient | HTTPクライアントのヘルパーコード群 |
|  |  | registry | DIコンテナフレームワークとコンテナで管理されるコンポーネントの定義およびHTTPリクエストルーティング定義 |
|  |  | logger | ロガー |

## Domainレイヤー

ここにはエンタープライズレベルのビジネスルールやオブジェクトが配置されている。  
今回のユースケースでは顧客の登録が行われるため、顧客と顧客に関するオブジェクトが定義されている。  

### /domain/customer/valueobjects.go

```go
type Customer interface {
  Name() string
  Email() string
}
```

Customerは生成時にのみ値が設定されるイミュータブルなオブジェクトであるValueObjectとして定義されている。  
プロジェクトのポリシーにもよるので一概には言えないが、安全で堅牢なコードという観点ではまずValueObject、どうしても必要な場合はEntityという順番でモデリングすることを推奨する。  
  
📔 **Note**  
Go言語では読み取り専用のプロパティを構造体に持たせることができないためInterfaceで表現している。  

### /domain/customer/builder.go

```go
type Builder interface {
  New(
    name string,
    email string,
  ) (Customer, error)
}
```

ここでは直接New関数を提供せず、BuilderインタフェースによるCustomerの生成を採用している。  
採用の理由はCustomerを使う他のレイヤーのコードがこのレイヤーの実装に依存しないでテストを可能とするためである。  

### /domain/customer/repository.go

```go
type Repository interface {
  Register(Customer) error
}
```

Customerの永続化を実現するRepositoryはこのレイヤーでInterfaceとして定義されている。  
前述の通り、これによりCustomerを操作したいものは外側のレイヤーの実装に影響されずに自身のコードに集中することができる。  
当然テストも雑音から解放され本来のロジックの検証にフォーカスできる。  

### /domain/creditscreening/domainservice.go

```go
type CreditScreeing interface {
  Perform(customer.Customer) error
}
```

このレイヤーには物や事を表現するValueObjectやEntity以外に、それらには当てはまらないビジネスロジックやイベントを表現するためのDomainServiceやDomainEventも配置される。  
ここでは顧客の審査を行うためのCreditscreeningがDomainServiceとして定義されている。  
これも実装ではなく仕様のみがこのレイヤーには配置される。理由は前述のRepositoryと全く同じである。  

## Usecaseレイヤー

ここに配置されるのは、名前の通りこのアプリケーションが実現するユースケースを表現するオブジェクトになる。  
今回は顧客の登録とそれに先立って顧客の審査を行うというユースケースが実現されることになる。  

### /usecase/customermanagement/spec.go

```go
type CustomerManagement interface {
  Register(customer.Customer) error
}
```

CustomerManagerは顧客を管理するメソッド群を提供するオブジェクトである。  
今回は顧客の登録というユースケースを実現するRegsiterメソッドのみ定義されている。  

### /usecase/customermanagement/impl.go

```go
import (
  "example.layeredarch/domain/creditscreening"
  "example.layeredarch/domain/customer"
)

type impl struct {
  customerRepository customer.Repository
  creditScreening    creditscreening.CreditScreeing
}
```

ここでは同じパッケージ内に実装コードも配置されている。  
実装と仕様(Interface)をわざわざ分けているのはこのオブジェクトを使う他のレイヤーがこのレイヤーの実装に依存しない様にするためである。  
そして、このオブジェクトが利用する他のレイヤーのオブジェクトは非公開プロパティとしてInterfaceの形で持つことにより(コードの`customerRepository`や`creditScreening`)このオブジェクトが他のレイヤーに依存することも防ぐことができる。  
常にレイヤー間の疎結合を保つことがてスタビリティを向上させ、テストのモチベーションを高めてくれる。  

📔 **Note**

```go
func New(
  customerRepository customer.Repository,
  creditScreening creditscreening.CreditScreeing,
) CustomerManagement {
  return &impl{
    customerRepository: customerRepository,
    creditScreening:    creditScreening,
  }
}
```

非公開プロパティの実装はこのオブジェクトが生成される時に外部から渡される。  
いわゆるDI(Dependency Injection)の中でも最もポピュラーなConstructor Injectionである。  
このサンプルプログラムでは**sarulabs/di**という軽量かつシンプルなDIコンテナによりそれを実現している。  

### /usecase/customermanagement/impl__register.go

```go
import (
  "example.layeredarch/domain/customer"
  "example.layeredarch/logger"
)

func (i *impl) Register(
  customer customer.Customer,
) error {
  if err := i.creditScreening.Perform(customer); err != nil {
    logger.Error(err)
    return err
  }

  if err := i.customerRepository.Register(customer); err != nil {
    logger.Error(err)
    return err
  }

  return nil
}
```

顧客を登録するというユースケースの本体がこのRegisterメソッドである。  
顧客を表すCustomerを受け取り、CreditScreeing#Performにより顧客の審査を行い、審査にパス(エラーが帰らない)したらRepository#Registerによって顧客を登録する。  
このメソッドが行うべきことは上記のとおり各外部オブジェクトを使って仕事を完了するのことのみであり、CreditScreeingやReositoryの実装がどのように審査や登録を行うかを意識する必要はないしすべきでもない。  
コードを見ればわかる通り、このメソッドに登場するオブジェクトは全てDomainレイヤー、つまりこのレイヤーより内側のレイヤーのもののみであり、外側のレイヤーの実装への依存は完全に排除されている。  
これこそがこのアーキテクチャーの醍醐味であり、絶対に犯すべきではないルールでもある。  

## Interfaceレイヤー

このレイヤーには永続化の実装や外部APIへのアクセスの実装のようなDomainレイヤーやUsecaseレイヤーに配置された仕様を実現するものと、WAFのような外界からのアクセスとの橋渡しを実現するものが配置される。  
  
### /repository/rdbcustomerrepository/impl.go

まずは永続化を担うRepositoryの実装から見てみよう。

```go
type impl struct {
  tx *sqlx.Tx
}
```

顧客の永続化を担当するcustomer.Repositoryの実装である。  
この実装は永続化層としてリレーショナルデータベースを対象とするため、リレーショナルデータベースを操作するためのトランザクションオブジェクト(`sqlx.Tx`)を非公開プロパティとしてもつ。  

```go
func New(tx *sqlx.Tx) customer.Repository {
  return &impl{tx: tx}
}
```

当然このオブジェクトもDIコンテナによりインジェクションされる。  
テストに際してはテスト用のトランザクションオブエクトを渡すことにより、実際のデータベースを対象に安全かつ健全な意味での恣意性を持ってテストを行うことができる。  

📔 **Note**

```go
import (
  "example.layeredarch/domain/customer"
  "github.com/jmoiron/sqlx"
)
```

上記のインポートを見て気づいただろうか。上のインポートはDomainレイヤーのものであるため依存の矢印は正しい向きを向いているが、下のインポートはどうだろうか。`jmoiron/sqlx`は外部ライブラリである。つまり所属するレイヤーはInfrastructureレイヤーになる。間違いなくこのレイヤーの外側のレイヤーである。  
このことを持ってClean Architecture(このドキュメントはClean Architectureを採用しているとは明言していないが、いまさらながらほぼそれであることは見ての通りである）は自家中毒だとか実現不可能だとかの評価をする向きもあることは事実である。ここで多くは語らないが(それをやるとそれだけでこれまで書いた量を軽く上回る)、Uncle Bobのブログをよく読めば、原理主義的な指摘が以下に的を射ていないかは理解できると思う。    

### /repository/rdbcustomerrepository/impl__register.go

```
func (i *impl) Register(customer customer.Customer) error {
  sql := `
INSERT INTO customers (
 name
,email
) values (
 :Name
,:Email
)
`
  qry, param, _ := sqlx.Named(
    sql,
    map[string]any{
      "Name":  customer.Name(),
      "Email": customer.Email(),
    },
  )

  if _, err := i.tx.Exec(qry, param...); err != nil {
    switch e := err.(type) {
    case *mysql.MySQLError:
      if e.Number == 1062 {
        logger.Error(e)
        return alreadyexistcustomer.New(customer)
      } else {
        logger.Error(e)
        return unexpected.New(err)
      }
    default:
      logger.Error(e)
      return unexpected.New(err)
    }
  } else {
    return nil
  }
}

```

顧客を登録するRegisterメソッドの実装である。  
実際の登録先はリレーショナルデータベースのテーブルであるため、ここでは公式が提供している`database/sql`の薄いラッパーライブラリである`jmoiron/sqlx`を使ってDBアクセスを実現している。  
見ての通り思い切り外部ライブラリのアイテムを使いまくりである。  
しかし、引数であるCustomerやエラー群のような内側のレイヤーとのやりとりに使われるものは全て内側のレイヤーに属するもののみである。  
これが内側のレイヤーを守りながら外側のレイヤーを利用しなければならないInterfaceレイヤーの役割そのものなのである。  
まさにインターフェースの名前が表すそのままの役所であると言えよう。  
  
📔 **Note**  
なぜリッチなORMではなく`jmoiron/sqlx`なのか。でもって、**sqlx**の控えめなORM的部分を使わずゴリゴリSQL描きたがるのか。  
ぶっちゃけてしまえば好きだから、というかリッチなORMが好きではないからである。  
一応好悪だけではない点も指摘しておくと、データベースの操作をほぼ隠蔽してしまうORMはCUDはさて置いたとしても、**R** つまりクエリーがどうしても大味になりがちなのがどうにも気持ちが悪いのである。あ、やはり好悪だ。世迷言と聞き流していただければ幸いである。  

### /gateway/httpcreditscreenig/impl.go

次は顧客信用審査を行うCreditScreeingの実装であるhttpcreditscreenigパッケージを見てみよう。  

```go
import (
  "example.layeredarch/domain/creditscreening"
  "example.layeredarch/gateway/httpcreditscreenig/translator"
  "example.layeredarch/httpclient"
)

type impl struct {
  client     httpclient.HttpClient
  translator translator.Translator
}

func New(
  client httpclient.HttpClient,
  translator translator.Translator,
) creditscreening.CreditScreeing {
  return &impl{client: client, translator: translator}
}
```

このパッケージはその名の通り顧客信用審査をHTTP接続によるREST-APIに委譲する。  
HTTP通信部分はInfrastructureレイヤーにあるhttpclientパッケージのHttpClientを使う。  
httpclientパッケージはHTTPクライアントライブラリである`dghubble/sling`を使ってHTTP通信を実現している。  

### /gateway/httpcreditscreenig/impl_perform.go

```go
func (i *impl) Perform(customer customer.Customer) error {
  data := i.translator.ToRequestData(customer)

  res, err := i.client.PostJson(data)
  if err != nil {
    logger.Error(err)
    return err
  }

  if res.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(res.Body)
    res := fmt.Sprintf(
      "StatusCode: %d, Body: %s",
      res.StatusCode,
      func() string {
        if body != nil {
          return string(body)
        } else {
          return ""
        }
      }(),
    )
    return unqualifiedcustomer.New(customer, res)
  }

  return nil
}
```

顧客の信用審査を行うPerformメソッドの実装である。  
やっていることは至って単純で、REST-APIを使って審査を行い結果がHTTP STATUSコード 200以外は審査に失敗した(審査にパスできなかった)と判断して、unqualifiedcustomerエラーを返すといった内容である。  
ここでは、DomainレイヤーのオブジェクトであるCustomerのREST-APIで使えるJSONフォーマットへの変換をになっているtranslatorパッケージに注目してみよう。  

### /gateway/httpcreditscreenig/translator/spec.go

```go
type Translator interface {
  ToRequestData(customer.Customer) CreditScreening
}
```

TranslatorインターフェースはCustomerをREST-APIが要求するJSONフォーマットを表す構造体CreditScreeningに変換するToRequestDataメソッドのみを持つ。  

### /gateway/httpcreditscreenig/translator/impl__torequestdata.go

```go
type CreditScreening struct {
  Name  string `json:"name"`
  Email string `json:"email"`
}
```

CreditScreeningはjsonへのマッピングタグを付与された構造体である。  
この構造体はREST-APIが要求するJSONフォーマットへの変換のためだけに存在し、他の用途に使われることはない。  

```go
func (i *impl) ToRequestData(customer customer.Customer) CreditScreening {
  return CreditScreening{
    customer.Name(),
    customer.Email(),
  }
}
```

変換を担当するメソッドの実装は至ってシンプルである。  
この例ではあまりにシンプルなため、Translatorなどと言う外付けのトランスレーターはともすれば冗長なだけに思える。  
しかし実プロダクトのプログラムでこの様に単純な変換だけしかないケースは稀であり、テスタビリティ、ソースの視認性と言う点だけとってもその価値は決して低くない。  

では、Interfaceレイヤー最後のパッケージであるrequesthandlerパッケージを見てみよう。  

### /requesthandler/spec.go

```go
import (
  "net/http"
)

type RequestHandler interface {
  HandleRequest(http.ResponseWriter, *http.Request) error
}
```

RequestHandlerインターフェースはその名通りHTTPリクエストをハンドリングするための唯一のメソッドであるHandleRequestを持つ。  
HTTPリクエストに応答するためのコードはこのインターフェースを実装することにより、InfrastructureレイヤーにあるWAFから呼び出される。  
以下では具体的な実装コードを見てその動きを理解することにしよう。  

### /requesthandler/customerpost/impl.go

```go
import (
  "example.layeredarch/requesthandler"
  "example.layeredarch/requesthandler/customerpost/translator"
  "example.layeredarch/usecase/customermanagement"
  "example.layeredarch/waf/defaulthandler"
)

type impl struct {
  *defaulthandler.DefaultHandler
  cutomerManagement customermanagement.CustomerManagement
  translator        translator.Translator
}
```

customerpostパッケージは/customerというパスに対して行われるPOSTリクエストに応答するRequestHandlerインターフェースの実装が配置される。  
上記の非公開プロパティを見ると、このパッケージの実装がDomainレイヤーで定義されているCustomerManagementインターフェースを使おうとしていることがわかる。加えてここでも専用のトランスレーター(**/requesthandler/customerpost/translator**)が用意されていることもわかる。  
上記の説明はRequestHandlerメソッドの実装に譲るとして、ここでは最後のプロパティである**defaulthandler.DefaultHandler**に注目してみよう。  

```go
type DefaultHandler struct{}
```

インポートを見るとわかる通り、この構造体はInfrastructureレイヤーのwafパッケージ配下にあるdefaulthandlerパッケージに定義されている。
この構造体はHTTPレスポンスの組み立てと書き出しを手助けするヘルパーメソッドの集合である。

| メソッド名 | 役割 |
| --- | --- |
| ResponseJson | JSONをボディにもつレスポンスを返す |
| ResponseError | エラーレスポンスを返す |
| ResponseEmpty | 空のレスポンスを返す |

📔 **Note**  
ご存知の通りGo言語に型の継承は無い。このため共通の処理を展開するにはEmbedded structsによるコンポジションに頼るしかない。  
しかない、とは言ったがComposition over inheritanceの概念はオブジェクト指向の割と初期の頃から言われていたことであり確かに有益な点が多い。  
Go言語の割り切りは個人的にはネガティブなのだが良いことも確かにある。少なくとも継承地獄は避けられる。  

### /requesthandler/customerpost/impl__handlerequest.go

```go
import (
  "net/http"

  "example.layeredarch/domain/errors/alreadyexistcustomer"
  "example.layeredarch/domain/errors/unqualifiedcustomer"
  "example.layeredarch/logger"
)

func (i *impl) HandleRequest(w http.ResponseWriter, r *http.Request) error {
  customer, err := i.translator.ToCustomer(r)
  if err != nil {
    logger.Error(err)
    i.ResponseError(http.StatusBadRequest, w, err)
    return err
  }

  if err := i.cutomerManagement.Register(customer); err != nil {
    switch e := err.(type) {
    case *unqualifiedcustomer.Error:
      logger.Error(e)
      i.ResponseError(http.StatusBadRequest, w, e)
    case *alreadyexistcustomer.Error:
      logger.Error(e)
      i.ResponseError(http.StatusConflict, w, e)
    default:
      logger.Error(e)
      i.ResponseError(http.StatusInternalServerError, w, e)
    }
    return err
  }

  logger.Info("Customer created. [name: %s, email: %s]", customer.Name(), customer.Email())
  i.ResponseEmpty(http.StatusCreated, w)
  return nil
}
```

HandleRequestの実装である。これはMVCにおけるControllerと位置付けとしては同じと考えて良い。  
なぜ  **”Contoroller”** ではなく **”Handler”** かと言えば、このアプリケーションはMVCでは全然無いからである。  
加えてGoの標準ライブラリであるhttpパッケージの文脈でもリクエストのハンドリングはHandlerとなっているためそちらに寄せたとも言える。  
長いものに巻かれると安心するオーディナリーピープルの性である。  

```go
  customer, err := i.translator.FromRequest(r)
  if err != nil {
    logger.Error(err)
    i.ResponseError(http.StatusBadRequest, w, err)
    return err
  }
```

ここも至ってシンプルな実装である。まずはトランスレータによってリクエストボディに入ってきたJSONをCustomerに変換している。  
変換に失敗した場合は不正なデータが検出されたと言うことなのでBadRequestを返す。  

```go
  if err := i.cutomerManagement.Register(customer); err != nil {
    switch e := err.(type) {
    case *unqualifiedcustomer.Error:
      logger.Error(e)
      i.ResponseError(http.StatusBadRequest, w, e)
    case *alreadyexistcustomer.Error:
      logger.Error(e)
      i.ResponseError(http.StatusConflict, w, e)
    default:
      logger.Error(e)
      i.ResponseError(http.StatusInternalServerError, w, e)
    }
    return err
  }
```

次にUsecaseレイヤーのcustomermanagementパッケージに処理を委譲してアプリケーションが期待するユースケースを実現している。  
委譲の結果なんらかの問題が発生した場合はその内容に合わせてエラーレスポンスを返している。  

| エラーオブジェクト | エラーの内容 | レスポンスコード |
| --- | --- | --- |
| unqualifiedcustomer.Error | 信用審査にパスできなかった | 400 BadRequest |
| unqualifiedcustomer.Error | 既に同じメールアドレスを持つ顧客が存在する | 409 Conflict |
| 上記以外の全てのエラー | 想定外のエラーが発生 | 500 InternalServerError |

```go
  logger.Info("Customer created. [name: %s, email: %s]", customer.Name(), customer.Email())
  i.ResponseEmpty(http.StatusCreated, w)
  return nil
```

ユースケースが期待通りの結果を得られたら **201 Created** を返して終わる。

📔 **Note**  
このメソッドはerrorを戻り値として要求している。  
HTTPリクエストの入り口であり出口であるこのメソッドがなぜerrorを返す必要があるのであろうか。  
素直に考えればレスポンスを返した時点で終わりであるし、errorのハンドリングを取りこぼしていることも無い。  
それでもerrorを返す理由はデータベーストランザクションの後始末に関係がある。  
詳しくはInfrastructureレイヤーの説明に譲るので興味がある向きはそちらを読んでいただきたい。  
ただしその内容は当ドキュメントの主題からは外れる内容であるため必読では無いことをお断りしておく。  

### /requesthandler/customerpost/translator/spec.go

トランスレータについては **/gateway/httpcreditscreenig/translator** でほぼ同じ内容を説明しているが、幾つかトピックがあるので軽く触れておく。  

```go
type Translator interface {
  ToCustomer(r *http.Request) (customer.Customer, error)
}
```

このトランスレーターはhttp.Requestオブジェクトを受け取ってリクエストボディの中身をCustomerに変換するメソッド **ToCustomer** のみを持つ。  

### /requesthandler/customerpost/translator/impl.go

```go
import "example.layeredarch/domain/customer"

type impl struct {
  customerBuilder customer.Builder
}
```

ここでは **customer.Builder** に注目したい。  
ToCustomerメソッドの中でCustomerを生成するのであるが、New関数による直接生成だとどうしてもDomainレイヤーの実装に依存することになってしまう。  
より内側のレイヤーなので依存に問題は無い様にも思えるが、できる限り仕様による依存に留めた方が結合度がより薄まる。テスタビリティも高まる。  
ただし幾分直感的では無くなることは否定できない。この辺りはトレードオフになるので開発方針でどちらにでも転ぶところではあるが、個人的には疎結合はメリットの方がデメリットより大きいのでできる限り追求したいところではある。  

### /requesthandler/customerpost/translator/impl__tocustomer.go

```go
func (i *impl) ToCustomer(
  r *http.Request,
) (
  customer.Customer,
  error,
) {
  var rc struct {
    Name  string `json:"name"`
    Email string `json:"email"`
  }
  if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
    return nil, err
  }

  return i.customerBuilder.New(rc.Name, rc.Email)
}
```

当ドキュメントの趣旨からは外れるが、ここでは匿名構造体(Anonumous struct)の使い方に言及しておく。  
匿名構造体は初期化が煩雑になる、再利用性が無いなどのデメリットもあるが、今回の様に一時的にしか使わない(リクエストボディのデコードだけ)場合有益である。他の要素にも共通して言えることだが、スコープはまず狭いところから始めることが堅牢なソフトウエアを作るコツである。  

## Infrastructureレイヤー

Infrastructureレイヤーは既に言及している通り、フレームワーク(例えばWAF)やドライバー類(データベースドライバー)などアプリケーション固有では無いものが配置される。改めて紹介すると、このサンプルソフトウエアでは以下のものが配置されている。  
  
| パッケージ | 内容 |
| --- | --- |
| waf | WEBアプリケーションフレームワークとそれをアプリケーションに適用するためのコード類 |
| httpclient | HTTPクライアントライブラリとヘルパーコード類 |
| registry | DIコンテナとコンテナの管理情報 |
| logger | ロガー |

レイヤードアーキテクチャーの具体的な使い方を説明すると言うこの章の趣旨に沿わないので各パッケージを詳細に説明することは控えるが、各レイヤーを裏側で繋ぐ役割を担う **registryパッケージ** と RequestHandlerが使われる **wafパッケージ** 配下の **routerパッケージ** について少し説明することにする。  

### /registry/registry.go

```go
import (
  "example.layeredarch/registry/defs"
  "github.com/sarulabs/di"
)

func New() (di.Container, error) {
  return build(
    defs.Server(),
    defs.RequestRouter(),
    defs.Database(),
    defs.Transaction(),
    defs.DomainCustomerBuilder(),
    defs.RdbCustomerReporitory(),
    defs.HttpCreditScreening(),
    defs.HttpCreditScreeningTranslator(),
    defs.RequestHandlerCustomerPost(),
    defs.RequestHandlerCustomerPostTranslator(),
    defs.UsecaseCustomerManagement(),
    defs.CreditScreeningHttpClient(),
  )
}
```

ここではDIコンテナの構築を行なっている。  
**registryパッケージ** 配下の **defsパッケージ** にはコンテナで管理したいものをコンテナに登録するために必要な情報が定義されている。  
各定義を取得するメソッドをパッケージスコープメソッドであるbuildの引数に列挙することによりコンテナへの登録が自動で行われる。  
以下では幾つかの定義情報を見ることにより、DIコンテナの動きについて軽く説明する。  
  
**defsパッケージ**には、コンテナに登録するコンポーネント毎にコンテナ要素の定義情報である **di.Defオブジェクト** を返す関数を記述したファイルを配置する。  

### /registry/defs/server.go

```go
func Server() di.Def {
  return di.Def{
    Name:  "server",
    Scope: di.App,
    Build: func(cnt di.Container) (any, error) {
      port := os.Getenv("LISTEN_PORT")
      router, err := cnt.SafeGet(RequestRouter().Name)
      if err != nil {
        return nil, err
      }
      return server.New(
        port,
        router.(waf.Router),
      ), nil
    },
  }
}
```

まずは、**wafパッケージ** の **Server** の定義を見てみよう。  
di.Defには以下の要素がある。  

| プロパティ | 内容 |
| --- | --- |
| Name | コンポーネントの名前。 他のコンポーネントにインジェクションするときなどにコンテナから取得するキーになる。 |
| Scope | その名の通りコンテナのスコープを表す。 スコープにはApp, Request, SubRequestの3種類がある。各スコープは先の記述順に包含関係にある。 |
| Build | コンテナに登録するコンポーネントを生成するメソッド。ここに記述されたコンポーネントの生成メソッドがコンテナにオブジェクトの取得依頼が来るたびに実行される。ただしAppスコープのコンポーネントはコンテナに対してシングルトンとなり、RequestスコープのコンポーネントはSubContainerの生成毎にシングルトンとなる |

上記のコードを見ると、ServerはAppスコープで、リッスンポートを環境変数から、Routerをコンテナから取得して生成時にインジェクションしていることがわかる。当サンプルアプリケーションではDIコンテナは一つしか構築しないため、AppスコープであるServerはアプリケーションに対してシングルトンのコンポーネントになることもわかる。  

### /registry/defs/requesthandler_customerpost.go

```go
func RequestHandlerCustomerPost() di.Def {
  return di.Def{
    Name:  "requesthandler/customerpost",
    Scope: di.Request,
    Build: func(cnt di.Container) (any, error) {
      uc, err := cnt.SafeGet(UsecaseCustomerManagement().Name)
      if err != nil {
        return nil, err
      }
      tl, err := cnt.SafeGet(RequestHandlerCustomerPostTranslator().Name)
      if err != nil {
        return nil, err
      }
      return customerpost.New(
        uc.(customermanagement.CustomerManagement),
        tl.(translator.Translator),
      ), nil
    },
  }
}
```

対して **customerpostパッケージ** の **HandlerRequest** の実装はRequestスコープで **usecaseパッケージ** 配下の **CustomerManagement** と自身が使うトランスレーターをコンテナから取得して生成時にインジェクションしている。このことからHandlerRequestはAppスコープのコンテナからSubContainerが生成されるたびに独立したインスタンスが生成されることがわかる。  
  
📔 **Notes**  
コンポーネントのスコープは以下の指針で決定することができる。  
  
- ステートレスなコンポーネントはAppスコープで問題ない。むしろメモリリソース節約のために積極的にAppスコープにするべきである。
- ステートフルなコンポーネントはそのライフサイクルによってスコープを選ぶ。例えばデータベーストランザクションをプロパティに持つRepositoryの実装コンポーネントはトランザクション境界に合わせてスコープを選択する必要がある。(当サンプルアプリケーションではHTTPリクエストがトランザクション境界なので後述するRouterの実装からもRequestスコープにする必要がある)
- コンテナがアプリケーションで唯一である場合、ステートレス, ステートフル別にかかわらずアプリケーションで唯一のインスタンスである必要がある場合はAppスコープを選択する

### /registry/router.go

**registryパッケージ** の RegisterRouterメソッドはWAFのリクエストルータにハンドラーを登録する役割を負う。    

```go
func RegisterRouter(router waf.Router) {
  router.Register(
    waf.HandlerDef{
      Method: httpmethod.Post,
      Path:   "/customer",
      Name:   defs.RequestHandlerCustomerPost().Name,
    },
  )
}
```

**registryパッケージ** の RegisterRouterメソッドはmain関数で以下の様に使われている。  

```go
  router, err := getRouter(reg)
  if err != nil {
    return nil
  }

  registry.RegisterRouter(router)
```

アプリケーションの起動時にコンテナからRouterを取得し、**RegisterRouter**にRouterを引き渡すことにより上記のコードの通りルータに **RequestHandler**が登録される。これにより、リクエストパスとHTTPメソッドおよびリクエストのハンドラーが関係づけられる。
リクエストパスとHTTPメソッドの組み合わせの数だけここでハンドラーとの関連付けを行うことにより複数のリクエストを処理可能となる。

📔 **Note**
  
**customerpostパッケージ** のHandlerRequest実装の[Note](https://www.notion.so/4637030949564628b5b4b5d6ebe73f55)で言及したデータベーストランザクションのハンドリングについて説明する。  

```go
    tx, err := i.beginTx(ctn)
    if err != nil {
      i.responseError(err, w)
      return
    }

    if err := handler.HandleRequest(w, r); err != nil {
      logger.Error(err)
      if e := tx.Rollback(); e != nil {
        logger.Error(e)
      }
    } else {
      if err := tx.Commit(); err != nil {
        logger.Error(err)
      }
    }
```

上記のコードは **wafパッケージ** 配下の **routerパッケージ** にある Routerインターフェースの実装コードの一部である。    
1行目の以下のコードでは、Routerコンポーネントがコンテナから生成する際にインジェクションされたコンテナ自身の参照からリクエストのハンドリング開始のタイミングで生成されたSabContainerからRequestスコープのデータベーストランザクションを取得している。    

```go
    tx, err := i.beginTx(ctn)
```

**sqlx.Tx**はコンテナから取得されると同時にトランザクションを開始する。  
そのトランザクションはRequestスコープ、すなわちHTTPリクエストがハンドリングされている間リクエスト毎に作られるSubContainerで管理され、データベーストランザクションが必要なコンポーネントにインジェクションされる。つまり、HTTPリクエストがトランザクション境界となるわけである。  
  
7行目ではリクエストハンドラーが呼び出されている。  

```go
    if err := handler.HandleRequest(w, r); err != nil {
```

例えば **customerpostパッケージ** で実装されているハンドラー内ではCustomerを扱うRepositoryがデータベースに顧客を永続化する箇所があるが、リクエストハンドラー内で何らかのエラー発生した場合9行目でトランザクションはロールバックする。  

```go
      if e := tx.Rollback(); e != nil {
```

リクエストハンドラーがエラーを起こすことなく終了すれば、13行目でトランザクションはコミットされる。  

```go
      if err := tx.Commit(); err != nil {
```

この様に、リクエストハンドラーがエラーを返す仕様になっていることによりトランザクションがHTTPリクエスト内で自動かつ適切に処理される。  
