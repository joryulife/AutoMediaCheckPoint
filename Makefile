#認証情報のjson設定用のコマンドです。

#jsonset: GCPの認証情報の入ったJSONを保存したpathに書き換えてください。
#    export GOOGLE_APPLICATION_CREDENTIALS = "~/AutoMediaCheckPoint/lib/~.json"

#外部ツールインストール
Setting:
    brew install annie
    brew install ffmpeg
    brew install mecab
    brew install mecab-ipadic
    git clone git@github.com:neologd/mecab-ipadic-neologd.git
    ./mecab-ipadic-neologd//bin/install-mecab-ipadic-neologd -n
    export CGO_LDFLAGS="-L/{libフォルダへのパス}/lib -lmecab -lstdc++"
    export CGO_CFLAGS="-I/{includeフォルダへのパス}/include"