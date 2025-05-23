import lume from "lume/mod.ts";
import googleFonts from "lume/plugins/google_fonts.ts"
import unocss from "lume/plugins/unocss.ts"
import nunjucks from "lume/plugins/nunjucks.ts"
import nav from "lume/plugins/nav.ts"
import lightningcss from "lume/plugins/lightningcss.ts"
import relativeUrls from "lume/plugins/relative_urls.ts"
import icons from "lume/plugins/icons.ts"
import inline from "lume/plugins/inline.ts"
import metas from "lume/plugins/metas.ts"

const site = lume({
  src: "./",
  dest: "build",
  server: {
    debugBar: false
  }
});

site.use(nunjucks({
  pageSubExtension: ""
}))
site.use(nav())
site.use(googleFonts({
  fonts: "https://fonts.googleapis.com/css2?family=Figtree:ital,wght@0,300..900;1,300..900&display=swap",
  cssFile: "styles.css"
}))
site.use(unocss({
  options: {
    theme: {
      fontFamily: {
        sans: "Figtree"
      },
      colors: {
        ui: {
          base: "rgb(15, 26, 21)",
          accent: "rgb(141, 169, 127)",
          text: "rgb(201, 216, 197)",
          subtext: "rgba(201, 216, 197, 0.7)",
          surface: "rgb(26, 41, 34)",
          elevated: "rgb(42, 59, 47)"
        }
      }
    },
  },
  cssFile: "styles.css"
}))
site.use(lightningcss())
site.use(relativeUrls())
site.use(icons())
site.use(inline())
site.use(metas())
site.add("https://unpkg.com/@unocss/reset@66.1.2/tailwind.css", "styles.css")
site.add("/assets")

site.data("year", new Date().getFullYear())

export default site;
