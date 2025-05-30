import lume from "lume/mod.ts";
import googleFonts from "lume/plugins/google_fonts.ts"
import unocss from "lume/plugins/unocss.ts"
import nunjucks from "lume/plugins/nunjucks.ts"
import nav from "lume/plugins/nav.ts"
import lightningcss from "lume/plugins/lightningcss.ts"
import icons from "lume/plugins/icons.ts"
import inline from "lume/plugins/inline.ts"
import metas from "lume/plugins/metas.ts"
import date from "lume/plugins/date.ts"
import codeHighlight from "lume/plugins/code_highlight.ts"

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
  fonts: "https://fonts.googleapis.com/css2?family=Figtree:ital,wght@0,300..900;1,300..900&family=Roboto+Mono:ital,wght@0,100..700;1,100..700&display=swap",
  cssFile: "styles.css",
  subsets: [
    "latin",
    "latin-ext"
  ]
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
    shortcuts: {
      'btn': `relative flex items-center gap-2 px-4 py-2 bg-ui-elevated text-ui-text font-medium text-sm rounded-md 
      shadow-[0_4px_0_rgb(26,41,34)] active:translate-y-1 active:shadow-none transition-all`,
      'btn-primary': `bg-ui-accent text-ui-base shadow-[0_4px_0_rgb(141,169,127,0.5)]`
    }
  },
  cssFile: "styles.css",
  reset: "tailwind"
}))
site.use(lightningcss())
site.use(icons())
site.use(inline())
site.use(metas())
site.use(date())
site.use(codeHighlight())

site.add("/assets/styles.css", "styles.css")
site.add("/assets")

// make an excerpt for each page to be used by the article preview
site.preprocess([".md"], (pages) => {
  for (const page of pages) {
    page.data.excerpt ??= (page.data.content as string).split(/<!--\s*more\s*-->/i)[0]
  }
})

// strip images from article previews
site.process([".html"], (pages) => {
  for (const page of pages) {
    const previewElements = page.document.querySelectorAll("article p")
    for (const previewElement of previewElements) {
      const images = previewElement.querySelectorAll("img")
      for (const image of images) {
        image.remove()
      }
    }
  }
})

site.data("year", new Date().getFullYear())

export default site;
