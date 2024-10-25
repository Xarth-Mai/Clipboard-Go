# Maintainer: Your Name <your_email@example.com>
pkgname=clipboard-go
pkgver=1.0.1
pkgrel=1
pkgdesc="Share clipboard between 💻Linux & 📱iOS"
arch=('any')
url="https://github.com/Xarth-Mai/Clipboard-Go"
license=('MPL2')
makedepends=('go')
depends=('xclip')
backup=(etc/clipboard-go/config.json)
source=("${pkgname}-${pkgver}.tar.gz::https://github.com/Xarth-Mai/Clipboard-Go/archive/refs/tags/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$srcdir/Clipboard-Go-${pkgver}"
  go get github.com/Xarth-Mai/EasyI18n-Go
  go build -o "${pkgname}" .
}

package() {
  install -Dm755 "${srcdir}/Clipboard-Go-${pkgver}/${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
  install -Dm644 "${srcdir}/Clipboard-Go-${pkgver}/config.json" "${pkgdir}/etc/clipboard-go/config.json"
  install -Dm644 "${srcdir}/Clipboard-Go-${pkgver}/clipboard-go.service" "${pkgdir}/usr/lib/systemd/user/${pkgname}.service"
}
