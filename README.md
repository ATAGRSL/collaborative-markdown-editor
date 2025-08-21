# 🎯 Collaborative Markdown Editor

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![WebSocket](https://img.shields.io/badge/WebSocket-Real--time-010101?style=for-the-badge&logo=websocket&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
[![HTML5](https://img.shields.io/badge/HTML5-5.0-E34F26?style=for-the-badge&logo=html5&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/HTML)
[![CSS3](https://img.shields.io/badge/CSS3-3.0-1572B6?style=for-the-badge&logo=css3&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/CSS)
[![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)

> **Gerçek zamanlı, çok kullanıcılı Markdown editörü** - Google Docs benzeri deneyim, ancak Markdown formatında! 🎉

Collaborative Markdown Editor, **Go** tabanlı gelişmiş bir gerçek zamanlı markdown editörüdür. WebSocket teknolojisi ile çoklu kullanıcıların aynı dokümanı eş zamanlı olarak düzenlemesine olanak tanır.

## 📋 **Repository Bilgileri**

- **GitHub Repository**: [https://github.com/ATAGRSL/collaborative-markdown-editor](https://github.com/ATAGRSL/collaborative-markdown-editor)
- **Demo**: [http://localhost:8080](http://localhost:8080) (lokal çalıştırma sonrası)
- **Lisans**: MIT License
- **Durum**: Aktif Geliştirme

## 🚀 Özellikler

### 🎨 **Çekirdek Özellikler**
- **Gerçek Zamanlı Senkronizasyon** ⚡: Bir kullanıcı yazdığında, değişiklikler anında diğer tüm kullanıcılara yansır
- **Oda Sistemi** 🏠: Benzersiz URL'ler ile farklı yazı odaları oluşturabilirsiniz
- **Anlık Markdown Önizlemesi** 👀: Solda Markdown metni, sağda HTML önizlemesi
- **WebSocket Tabanlı** 🔌: Go'nun goroutine'leri ile verimli gerçek zamanlı iletişim
- **Responsive Tasarım** 📱: Modern ve kullanıcı dostu arayüz

### 👥 **Kullanıcı Yönetimi**
- **Active Users Listesi** 👥: Odaya katılan tüm kullanıcıları gerçek zamanlı görüntüleme
- **Renk Kodlaması** 🎨: Her kullanıcı için benzersiz renkler ve avatar'lar
- **Otomatik Kullanıcı Adları** 🔤: Rastgele oluşturulan kullanıcı adları
- **Gerçek Zamanlı Katılma/Çıkma** 📈: Kullanıcıların giriş/çıkışlarını anlık takip

### 💻 **Gelişmiş Editör Özellikleri**
- **Markdown Syntax Highlighting** 🌈: Kod blokları için syntax highlighting
- **Auto-save** 💾: Yazılar otomatik olarak kaydedilir
- **Cursor Position Tracking** 📍: Diğer kullanıcıların imleç konumlarını görme
- **Live Preview** 🔄: Markdown'dan HTML'e anlık dönüşüm
- **Keyboard Shortcuts** ⌨️: Markdown yazımı için klavye kısayolları



## Kurulum ve Çalıştırma

### 💻 Gereksinimler
- **Go 1.21+** kurulu olmalı
- **Web tarayıcı** (Chrome, Firefox, Safari, Edge)
- **Internet bağlantısı** (CDN kütüphaneleri için)

### 🚀 Kurulum Adımları

1. **Proje dizinine gidin:**
   ```bash
   cd collaborative-markdown-editor
   ```

2. **Bağımlılıkları yükleyin:**
   ```bash
   go mod tidy
   ```

3. **Server'ı başlatın:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Web tarayıcıda açın:**
   ```
   http://localhost:8080
   ```

### 🧪 Hızlı Test

Sunucu başladıktan sonra tarayıcıda `http://localhost:8080` adresini açın ve:
- ✅ Ana sayfa yüklendi mi?
- ✅ "✨ Create New Room" butonu çalışıyor mu?
- ✅ Yeni oda oluşturulabiliyor mu?
- ✅ Markdown editörü çalışıyor mu?
- ✅ Active Users listesi görünüyor mu?

### 🎮 **Demo Kullanım**

#### **1. İlk Kullanıcı - Oda Oluşturma**
1. `http://localhost:8080` adresini açın
2. **"✨ Create New Room"** butonuna tıklayın
3. URL'de `/room/abc123` gibi bir oda ID'si göreceksiniz
4. Active Users listesinde kendinizi göreceksiniz

#### **2. İkinci Kullanıcı - Odaya Katılma**
1. Yeni bir tarayıcı sekmesi açın
2. `http://localhost:8080` adresine gidin
3. Room ID kutusuna `abc123` yazın
4. **"Join Room"** butonuna tıklayın
5. Her iki sekmedeki Active Users listesinde 2 kullanıcı göreceksiniz

#### **3. Eş Zamanlı Düzenleme**
1. Bir sekmede markdown yazın: `# Merhaba Dünya`
2. Diğer sekmedeki değişiklikleri anlık görün
3. Her iki sekmedeki preview panelini kontrol edin

## 📖 Kullanım

### 🏠 **Oda Yönetimi**

#### **Yeni Oda Oluşturma**
1. Ana sayfaya gidin (`http://localhost:8080`)
2. **"✨ Create New Room"** butonuna tıklayın
3. Otomatik olarak yeni bir odaya yönlendirileceksiniz
4. URL'de benzersiz oda ID'si görünecek

#### **Mevcut Odaya Katılma**
1. Ana sayfaya gidin
2. **Room ID** kutusuna katılmak istediğiniz oda ID'sini girin
3. **"Join Room"** butonuna tıklayın
4. Aynı odaya katılmış olacaksınız

### ✍️ **Markdown Editörü Kullanımı**

#### **Temel Özellikler**
- **Soldaki Panel**: Markdown metni yazın
- **Sağdaki Panel**: Anlık HTML önizlemesi
- **Gerçek Zamanlı**: Tüm değişiklikler diğer kullanıcılara senkronize edilir

#### **Kullanılabilir Markdown Özellikleri**

````markdown
# Başlık 1
## Başlık 2

**Kalın metin**
*İtalik metin*

- Liste öğesi 1
- Liste öğesi 2

[Bağlantı](https://example.com)

`inline kod`

```js
// Kod bloğu
console.log('Merhaba!');
```

> Alıntı metni
````

#### **Klavye Kısayolları**
- **Ctrl/Cmd + B**: Kalın metin
- **Ctrl/Cmd + I**: İtalik metin
- **Ctrl/Cmd + K**: Bağlantı ekleme
- **Enter**: Yeni satır

### 👥 **Çoklu Kullanıcı Özellikleri**

#### **Active Users Listesi**
- Sağ panelde **"👥 Active Users"** bölümünde tüm aktif kullanıcıları görürsünüz
- Her kullanıcı benzersiz bir renkle temsil edilir
- Kullanıcı adları otomatik olarak oluşturulur
- Gerçek zamanlı giriş/çıkış takibi

#### **Eş Zamanlı Düzenleme**
- Birden fazla kullanıcı aynı anda yazabilir
- Değişiklikler anında tüm kullanıcılara yansır
- Cursor pozisyonları takip edilir
- Çakışma durumunda otomatik çözümleme

## 🏗️ Proje Mimarisi

### 📂 **Dosya Yapısı**

```
collaborative-markdown-editor/
├── cmd/server/main.go          # Ana HTTP/WebSocket sunucusu
├── internal/
│   ├── client/client.go        # WebSocket client yönetimi
│   ├── hub/hub.go             # Oda ve kullanıcı hub'ı
│   ├── ot/
│   │   ├── manager.go         # Operational Transformation
│   │   └── operation.go       # OT operasyonları
│   └── user/user.go           # Kullanıcı yönetimi
├── go.mod                     # Go module tanımı
├── go.sum                     # Bağımlılık checksum'ları
└── README.md                  # Bu dosya
```

### 🏛️ **Ana Bileşenler**

#### **Backend Altyapısı**
- **Hub**: Tüm odaları ve bu odalardaki client'ları yönetir
- **Client**: Her WebSocket bağlantısını temsil eder
- **User Manager**: Kullanıcı kayıt ve yönetimi
- **OT Manager**: Operational Transformation ile çakışma çözümü
- **Goroutines**: Her client için readPump ve writePump goroutineleri

#### **Frontend Altyapısı**
- **HTML5**: Modern ve semantic HTML yapısı
- **CSS3**: Responsive tasarım ve animasyonlar
- **JavaScript ES6+**: Modern JavaScript özellikleri
- **WebSocket API**: Gerçek zamanlı iletişim
- **Marked.js**: Markdown parsing (CDN)

### 🔌 **API Endpoints**

| Endpoint | Method | Açıklama |
|----------|--------|----------|
| `GET /` | GET | Ana sayfa ve oda oluşturma arayüzü |
| `GET /room/{roomId}` | GET | Belirli bir odanın editör sayfası |
| WebSocket `/ws/{roomId}` | WebSocket | Gerçek zamanlı mesajlaşma endpoint'i |

### 📊 **Veri Akışı**

```
Client Request → HTTP Server → Hub → Room → Clients
                      ↓
WebSocket Connection → Message Processing → OT Manager
                      ↓
User List Updates → Real-time Broadcast → All Clients
```

## 🔧 Geliştirme

### 🚀 **Yeni Özellikler Eklemek**

#### **Backend Geliştirme**
1. **Hub'a yeni metodlar ekleyin** (`internal/hub/hub.go`)
2. **Client davranışını genişletin** (`internal/client/client.go`)
3. **HTTP endpoint'leri ekleyin** (`cmd/server/main.go`)
4. **Kullanıcı yönetimi özelliklerini geliştirin** (`internal/user/user.go`)

#### **Frontend Geliştirme**
1. **HTML template'lerini güncelleyin** (main.go içindeki template'ler)
2. **CSS stillerini genişletin**
3. **JavaScript işlevselliğini artırın**
4. **Yeni UI bileşenleri ekleyin**

### 🧪 **Test ve Debug**

#### **Çoklu Kullanıcı Testi**
Çoklu tarayıcı sekmesi açarak gerçek zamanlı senkronizasyonu test edebilirsiniz:
1. İki farklı tarayıcı sekmesi açın
2. Aynı odaya katılın
3. Bir sekmede yazın ve diğer sekmedeki değişiklikleri gözlemleyin
4. Active Users listesinin doğru çalıştığını kontrol edin

#### **Debugging İpuçları**
```bash
# Server loglarını izleme
go run cmd/server/main.go

# Port kullanımını kontrol etme
lsof -i :8080

# WebSocket bağlantılarını izleme
# Tarayıcı DevTools → Network → WS tab
```

## 🛠️ Sorun Giderme

### 🚨 **Sık Görülen Sorunlar**

#### **Port 8080 Kullanımda**
```bash
# Hangi process kullanıyor?
lsof -i :8080

# Process'i durdur
kill -9 <PID>

# Alternatif port kullan
go run cmd/server/main.go # Port 8080'de çalışır
```

#### **Bağımlılık Sorunları**
```bash
# Bağımlılıkları yeniden yükle
go mod tidy

# Cache temizle
go clean -modcache
```

#### **WebSocket Bağlantı Sorunları**
- Firewall ayarlarını kontrol edin
- Tarayıcıda WebSocket desteğini kontrol edin
- Server loglarında bağlantı hatalarını inceleyin

#### **Markdown Parse Hataları**
- İnternet bağlantısını kontrol edin (CDN kütüphaneleri için)
- Tarayıcı console'da JavaScript hatalarını kontrol edin
- Marked.js versiyon uyumluluğunu kontrol edin

### 📞 **Destek**

Herhangi bir sorun yaşarsanız:
1. **Server loglarını kontrol edin** (`go run cmd/server/main.go` çıktısı)
2. **Tarayıcı console'unda hataları kontrol edin** (F12 → Console)
3. **Network tab'ında WebSocket bağlantısını kontrol edin**
4. **Port çakışması olup olmadığını kontrol edin**

## 📊 Proje İstatistikleri

### 📈 **Kod Metrikleri**
- **Toplam Go Dosyası**: 6 dosya
- **Toplam Kod Satırı**: ~2,000+ satır (backend + frontend)
- **Ana Paketler**: client, hub, ot, user
- **WebSocket Endpoint'leri**: 1 adet
- **HTTP Route'ları**: 2 adet

### 🎯 **Özellik Sayıları**
- **Çekirdek Özellikler**: 5 temel özellik
- **Kullanıcı Yönetimi**: 4 kullanıcı özelliği
- **Editör Özellikleri**: 5 editör özelliği
- **API Endpoint'leri**: 3 endpoint
- **Markdown Özellikleri**: 10+ syntax desteği

### 📊 **Performans Metrikleri**
- **Gerçek Zamanlı Bağlantı**: WebSocket tabanlı
- **Eş Zamanlı Kullanıcı**: Sınırsız (sistem kaynaklarına bağlı)
- **Oda Kapasitesi**: Her oda için sınırsız kullanıcı
- **Yanıt Süresi**: < 100ms (lokal)
- **Bellek Kullanımı**: ~50MB (idle durumda)

### 🏗️ **Teknik Altyapı**
- **Backend**: Go 1.21+ (net/http, goroutines)
- **WebSocket**: Gorilla WebSocket
- **Frontend**: HTML5, CSS3, JavaScript ES6+
- **Markdown Parser**: Marked.js (CDN)
- **Database**: In-memory storage

## 🏆 Katkıda Bulunma

### 🤝 **Nasıl Katkıda Bulunabilirsiniz?**

1. **Fork edin**: Bu repository'yi fork edin
2. **Feature branch oluşturun**:
   ```bash
   git checkout -b feature/yeni-ozellik
   ```
3. **Değişikliklerinizi yapın**
4. **Commit edin**:
   ```bash
   git commit -am 'Yeni özellik eklendi'
   ```
5. **Push edin**:
   ```bash
   git push origin feature/yeni-ozellik
   ```
6. **Pull Request oluşturun**

### 🎯 **İleri Düzey Geliştirme Önerileri**

#### **Özellik İyileştirmeleri**
- [ ] Kullanıcı kimlik doğrulama sistemi
- [ ] Dosya kaydetme/yükleme sistemi
- [ ] Markdown export (PDF, DOCX)
- [ ] Tema değiştirme özelliği
- [ ] Syntax highlighting iyileştirmesi

#### **Performans İyileştirmeleri**
- [ ] Database entegrasyonu (SQLite/PostgreSQL)
- [ ] Redis ile session yönetimi
- [ ] CDN entegrasyonu
- [ ] Load balancing desteği

#### **Güvenlik İyileştirmeleri**
- [ ] Rate limiting
- [ ] Input sanitization
- [ ] XSS protection
- [ ] CSRF protection

## 🏆 Lisans

Bu proje **eğitim amaçlı** oluşturulmuştur ve **MIT lisansı** altında lisanslanmıştır.

## 🙏 Teşekkürler

Bu projeyi kullanmanız ve ilgilenmeniz için teşekkür ederim! Proje geliştirme sürecinde kullandığım teknolojilere ve kaynaklara minnettarım:

### 🛠️ **Kullanılan Teknolojiler**
- **Go**: Backend geliştirme ve concurrency
- **Gorilla WebSocket**: WebSocket server implementation
- **Marked.js**: Markdown parsing ve rendering
- **HTML5/CSS3**: Modern web teknolojileri
- **JavaScript ES6+**: Client-side functionality

### 📚 **Kaynaklar**
- **Go Documentation**: Backend geliştirme rehberi
- **MDN Web Docs**: Web API referansları
- **WebSocket RFC**: Protokol spesifikasyonu
- **Markdown Guide**: Syntax referansı

---

## 📅 **Son Güncellemeler**

### **v1.0.0** - İlk Yayın (Current)
- ✅ Gerçek zamanlı çok kullanıcılı markdown editörü
- ✅ WebSocket tabanlı iletişim sistemi
- ✅ Active Users listesi ile kullanıcı takibi
- ✅ Responsive tasarım
- ✅ Operational Transformation desteği
- ✅ Markdown live preview
- ✅ Oda tabanlı sistem



---

**🎉 Collaborative Markdown Editor** - Modern, gerçek zamanlı markdown editörü! 🚀

*Herhangi bir sorun yaşarsanız veya önerileriniz varsa [GitHub Issues](https://github.com/ATAGRSL/collaborative-markdown-editor/issues) açmaktan çekinmeyin!*
