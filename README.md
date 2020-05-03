# Preview Imager

Generate Image Previews for Open Graph tags.

## Example

![Image of an example][image-example]

## Usage
Install any ttf font (`npm install --save @ibm/plex`).

```js
const { PreviewImager } = require("@rensatsu/preview-imager");
const path = require("path");

const fontPath = path.join(
    __dirname,
    "/node_modules/@ibm/plex/IBM-Plex-Sans/fonts/complete/ttf/"
);

(async () => {
    const pri = new PreviewImager("node-out.png", "test");
    pri.set("siteName", "Test Blog");
    pri.set("colorBg", "#665eba");
    pri.set("colorBgOverlay", "#665eba9f");
    pri.set("colorBgOverlay", "#0000009f");
    pri.set("paddingX", "50");
    pri.set("paddingY", "50");
    pri.set("lineSpacing", "1.75");
    pri.set("imagePath", path.join(__dirname, "/wallpaper.jpg"));
    pri.set("logoPath", path.join(__dirname, "/logo-32.png"));
    pri.set("logoX", "50");
    pri.set("logoY", "50");
    pri.set("paddingYSite", "25");
    pri.set("blurStrength", "5");
    pri.set("fontTitle", path.join(fontPath, "/IBMPlexSans-Bold.ttf"));
    pri.set("fontSiteName", path.join(fontPath, "/IBMPlexSans-Regular.ttf"));

    await pri.generate().catch(e => {
        console.error(e, e.output || null);
    });
})();
```

[image-example]: https://cdn.jsdelivr.net/gh/rensatsu/preview-imager@latest/.repository/screenshot-1.png
