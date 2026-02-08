package components

import (
	"fmt"

	"github.com/Mau005/KraynoSerer/configuration"
)

type Components struct{}

func (c *Components) CreateForm(action, method, content string) string {
	return fmt.Sprintf(`
	<div>
		<form action="%s" method="%s">
			%s
			<button type="submit">ntivar</button>
		</form>
	<div>
	`, action, method, content)
}
func (c *Components) CreateFormImput(typeInput, idText, title, value string, required bool) string {
	content := `<input type="%s" id="%s" name="%s" value="%s" required>`
	if required {
		content = `<input type="%s" id="%s" name="%s" value="%s">`
	}

	return fmt.Sprintf(`
	<label for="%s">%s:</label>
	%s
	`, idText, title, fmt.Sprintf(content, typeInput, idText, idText, value))
}

func (c *Components) CreateFormTextArea(idTextArea, title, value string) string {
	return fmt.Sprintf(`
	<label for="%s">%s:</label>
	<textarea id="%s" name="%s" required>%s</textarea>
	`, idTextArea, title, idTextArea, idTextArea, value)
}
func (c *Components) CreateFormButton(nameButton string) string {
	return fmt.Sprintf(`
	<button type="submit">%s</button>
	`, nameButton)
}

func (c *Components) CreateLabelADiv(content string) string {
	return c.CreateDiv(fmt.Sprintf("<a>%s</a>", content))
}

func (c *Components) CreateLabelA(content string) string {
	return fmt.Sprintf("<a>%s</a>", content)
}

func (c *Components) CreateDiv(content string) string {
	return fmt.Sprintf("<div>%s</div>", content)
}

func (c *Components) CreateButtonForm(method, url, nameButton string) string {

	book := `
	<form action="%s" method="%s">
		<button type="submit" class="btn btn-primary"> %s </button>
	</form>
	`
	return fmt.Sprintf(book, url, method, nameButton)

}

func (c *Components) CreateTable(content string) string {

	return fmt.Sprintf(`
	<table class="table">	
		%s
	<table>
	`, content)
}

func (c *Components) CreateColsTable(attrCol ...string) string {

	resultValues := ""
	for _, value := range attrCol {
		resultValues += fmt.Sprintf(`<th scope="col">%s</th>`, value)
	}

	result := fmt.Sprintf(`
	<tr>
		%s
    </tr>
	`, resultValues)

	return fmt.Sprintf(`
	<thead class="thead-dark">
	%s
  	</thead>
	`, result)
}

func (c *Components) CreateRowsTable(attrRows ...string) string {
	procesing := ""
	for _, values := range attrRows {
		procesing += fmt.Sprintf(` <td>%s</td>`, values)
	}

	return fmt.Sprintf(`
	<tr>
	%s
	</tr>
	`, procesing)
}

func (c *Components) CreateRowsTableFinally(content string) string {
	return fmt.Sprintf(`
	<tbody>
	%s
	</tbody>
	`, content)
}

func (c *Components) CrceateTitle(content string) string {

	return fmt.Sprintf("<h1>%s</h1>", content)
}

func (c *Components) CreateTitle(nameTitle string) string {
	return configuration.NameDefaultServer + " " + nameTitle
}

func (c *Components) CreateMetaDefault(title, description string) string {
	if description == "" {
		description = "Social network for Tibia, tools, share statistics, share photos and videos of your exploits in this beautiful game"
	}
	return fmt.Sprintf(`
        <meta name="description" content="%s">
        <meta name="author" content="AinhoSoft">
        <meta property="og:title" content="%s">
        <meta property="og:url" content="/">
        <meta property="og:description" content="%s">
        <meta name="author" content="AinhoSoft">
	`, description, c.CreateTitle(title), description)
}

func (c *Components) CreateLink() string {
	return `
 		<link rel="icon" type="image/png" href="/static/assets/images/favicon.png">
        <!-- START: Styles -->
        <!-- Google Fonts -->
        <link href="https://fonts.googleapis.com/css?family=Roboto+Condensed:300i,400,400i,700%7cMarcellus+SC" rel="stylesheet">
        <!-- Bootstrap -->
        <link rel="stylesheet" href="/static/assets/vendor/bootstrap/dist/css/bootstrap.min.css">
        <!-- FontAwesome -->
        <script defer src="/static/assets/vendor/fontawesome-free/js/all.js"></script>
        <script defer src="/static/assets/vendor/fontawesome-free/js/v4-shims.js"></script>
        <!-- IonIcons -->
        <link rel="stylesheet" href="/static/assets/vendor/ionicons/css/ionicons.min.css">
        <!-- Revolution Slider -->
        <link rel="stylesheet" href="/static/assets/vendor/revolution/css/settings.css">
        <link rel="stylesheet" href="/static/assets/vendor/revolution/css/layers.css">
        <link rel="stylesheet" href="/static/assets/vendor/revolution/css/navigation.css">
        <!-- Flickity -->
        <link rel="stylesheet" href="/static/assets/vendor/flickity/dist/flickity.min.css">
        <!-- Photoswipe -->
        <link rel="stylesheet" href="/static/assets/vendor/photoswipe/dist/photoswipe.css">
        <link rel="stylesheet" href="/static/assets/vendor/photoswipe/dist/default-skin/default-skin.css">
        <!-- DateTimePicker -->
        <link rel="stylesheet" href="/static/assets/vendor/jquery-datetimepicker/build/jquery.datetimepicker.min.css">
        <!-- Summernote -->
        <link rel="stylesheet" href="/static/assets/vendor/summernote/dist/summernote-bs4.css">
        <!-- GODLIKE -->
        <link rel="stylesheet" href="/static/assets/css/godlike.css">
        <!-- Custom Styles -->
        <link rel="stylesheet" href="/static/assets/css/custom.css">
        <!-- END: Styles -->
        <!-- jQuery -->
        <script src="/static/assets/vendor/jquery/dist/jquery.min.js"></script>
	
	`
}

func (c *Components) CreatePreload() string {
	return `
        <div class="nk-preloader">
            <div class="nk-preloader-bg" style="background-color: #000;" data-close-frames="23" data-close-speed="1.2" data-close-sprites="/static/assets/images/preloader-bg.png" data-open-frames="23" data-open-speed="1.2" data-open-sprites="./assets/images/preloader-bg-bw.png">
            </div>
            <div class="nk-preloader-content">
                <div>
                    <img class="nk-img" src="/static/assets/images/logo.png" alt="TibiaKray V1.0.0" width="170">
                    <div class="nk-preloader-animation"></div>
                </div>
            </div>
            <div class="nk-preloader-skip">Skip</div>
        </div>

        <div class="nk-page-background op-8" style="background-image: url('/static/assets/images-svg/background.png');"></div>
        
        <div class="nk-page-background-audio d-none" data-audio="/static/assets/mp3/purpleplanetmusic-desolation.mp3" data-audio-volume="30" data-audio-autoplay="true" data-audio-loop="true" data-audio-pause-on-page-leave="true"></div>

		<div class="nk-page-border">
            <div class="nk-page-border-t"></div>
            <div class="nk-page-border-r"></div>
            <div class="nk-page-border-b"></div>
            <div class="nk-page-border-l"></div>
        </div>

		<header class="nk-header nk-header-opaque">
	`
}

func (c *Components) CreateFooter(lenguaje map[string]string) string {
	return fmt.Sprintf(`
	        <div class="nk-gap-4"></div>
            <div class="nk-gap-4"></div>
            <div class="nk-gap-4"></div>
            <footer class="nk-footer nk-footer-parallax nk-footer-parallax-opacity">
                <img class="nk-footer-top-corner" src="/static/assets/images/footer-corner.png" alt="">
                <div class="container">
                    <div class="nk-gap-2"></div>
                    <div class="nk-gap"></div>
                    <p> &copy; 2023 AinhoSoft.</p>
                    <div class="nk-footer-links">
                        <a href="#" class="link-effect">%s</a> <span>|</span> <a href="#" class="link-effect">%s</a>
                    </div>
                    <div class="nk-gap-4"></div>
                </div>
            </footer>
	`, lenguaje["termservice"], lenguaje["privacypolicy"])
}

func (c *Components) CreateButtonVolumen() string {
	return `
	<div class="nk-side-buttons nk-side-buttons-visible">
            <ul>
                <li>
                    <span class="nk-btn nk-btn-lg nk-btn-icon nk-bg-audio-toggle">
                        <span class="icon">
                            <span class="ion-android-volume-up nk-bg-audio-pause-icon"></span>
                            <span class="ion-android-volume-off nk-bg-audio-play-icon"></span>
                        </span>
                    </span>
                </li>
                <li class="nk-scroll-top">
                    <span class="nk-btn nk-btn-lg nk-btn-icon">
                        <span class="icon ion-ios-arrow-up"></span>
                    </span>
                </li>
            </ul>
        </div>
	`
}

func (c *Components) CreateLogin(lenguaje map[string]string) string {
	return fmt.Sprintf(`
	        <div class="nk-sign-form">
            <div class="nk-gap-5"></div>
            <div class="container">
                <div class="row">
                    <div class="col-lg-4 offset-lg-4 col-md-6 offset-md-3">
                        <div class="nk-sign-form-container">
                            <div class="nk-sign-form-toggle h3">
                                <a href="#" class="nk-sign-form-login-toggle active">%s</a>
                                <a href="#" class="nk-sign-form-register-toggle">%s</a>
                            </div>
                            <div class="nk-gap-2"></div>
                            <!-- START: Login Form -->
                            <form class="nk-sign-form-login active" method="POST" action="/login">
                                <input class="form-control" type="text" name="user" id="user" placeholder="%s / Email">
                                <div class="nk-gap-2"></div>
                                <input class="form-control" type="password" name="passworduser" id="passworduser" placeholder="%s">
                                <div class="nk-gap-2"></div>
                                <div class="form-check float-left">
                                    <label class="form-check-label">
                                        <input type="checkbox" class="form-check-input"> %s </label>
                                </div>
                                <button class="nk-btn nk-btn-color-white link-effect-4 float-right">%s</button>
                                <div class="clearfix"></div>
                                <div class="nk-gap-1"></div>
                                <a class="nk-sign-form-lost-toggle float-right" href="#">%s?</a>
                            </form>
                            <!-- END: Login Form -->
                            <!-- START: Lost Password Form -->
                            <form class="nk-sign-form-lost" action="/recovery" method="POST">
                                <input class="form-control" type="text" name="recovery" placeholder="%s / Email">
                                <div class="nk-gap-2"></div>
                                <button class="nk-btn nk-btn-color-white link-effect-4 float-right">%s</button>
                            </form>
                            <!-- END: Lost Password Form -->
                            <!-- START: Register Form -->
                            <form class="nk-sign-form-register" action="/create_user" method="POST">
                                <input class="form-control" type="text" name="username" id="username" placeholder="%s">
                                <div class="nk-gap-2"></div>
                                <input class="form-control" type="password" name="password" id="password" placeholder="%s">
                                <div class="nk-gap-2"></div>
                                <input class="form-control" type="password"  name="passwordTwo" id="passwordTwo"  placeholder="%s %s">
                                <div class="nk-gap-2"></div>
                                <input class="form-control" type="email" name="email" id="email" placeholder="Email">
                                <div class="nk-gap-2"></div>
                                <div class="form-check float-left">
                                    <label class="form-check-label">
                                        <input type="checkbox" name="policy" id="policy" class="form-check-input"> %s </label>
                                </div>
                                <button class="nk-btn nk-btn-color-white link-effect-4 float-right">%s</button>
                            </form>
                            <!-- END: Register Form -->
                        </div>
                    </div>
                </div>
            </div>
            <div class="nk-gap-5"></div>
        </div>
	`, lenguaje["login"], lenguaje["register"], lenguaje["user"], lenguaje["password"], lenguaje["remember"], lenguaje["login"], lenguaje["lostpassword"], lenguaje["user"], lenguaje["getin"],
		lenguaje["user"], lenguaje["password"], lenguaje["repeat"], lenguaje["password"], lenguaje["acceptpolicies"], lenguaje["register"])

}

func (c *Components) CreateScripts() string {
	return `
	  	<script src="/static/assets/vendor/object-fit-images/dist/ofi.min.js"></script>
        <!-- GSAP -->
        <script src="/static/assets/vendor/gsap/dist/gsap.min.js"></script>
        <script src="/static/assets/vendor/gsap/dist/ScrollToPlugin.min.js"></script>
        <!-- Popper -->
        <script src="/static/assets/vendor/popper.js/dist/umd/popper.min.js"></script>
        <!-- Bootstrap -->
        <script src="/static/assets/vendor/bootstrap/dist/js/bootstrap.min.js"></script>
        <!-- Sticky Kit -->
        <script src="/static/assets/vendor/sticky-kit/dist/sticky-kit.min.js"></script>
        <!-- Jarallax -->
        <script src="/static/assets/vendor/jarallax/dist/jarallax.min.js"></script>
        <script src="/static/assets/vendor/jarallax/dist/jarallax-video.min.js"></script>
        <!-- imagesLoaded -->
        <script src="/static/assets/vendor/imagesloaded/imagesloaded.pkgd.min.js"></script>
        <!-- Flickity -->
        <script src="/static/assets/vendor/flickity/dist/flickity.pkgd.min.js"></script>
        <!-- Isotope -->
        <script src="/static/assets/vendor/isotope-layout/dist/isotope.pkgd.min.js"></script>
        <!-- Photoswipe -->
        <script src="/static/assets/vendor/photoswipe/dist/photoswipe.min.js"></script>
        <script src="/static/assets/vendor/photoswipe/dist/photoswipe-ui-default.min.js"></script>
        <!-- Typed.js -->
        <script src="/static/assets/vendor/typed.js/lib/typed.min.js"></script>
        <!-- Jquery Validation -->
        <script src="/static/assets/vendor/jquery-validation/dist/jquery.validate.min.js"></script>
        <!-- Jquery Countdown + Moment -->
        <script src="/static/assets/vendor/jquery-countdown/dist/jquery.countdown.min.js"></script>
        <script src="/static/assets/vendor/moment/min/moment.min.js"></script>
        <script src="/static/assets/vendor/moment-timezone/builds/moment-timezone-with-data.min.js"></script>
        <!-- Hammer.js -->
        <script src="/static/assets/vendor/hammerjs/hammer.min.js"></script>
        <!-- NanoSroller -->
        <script src="/static/assets/vendor/nanoscroller/bin/javascripts/jquery.nanoscroller.js"></script>
        <!-- SoundManager2 -->
        <script src="/static/assets/vendor/soundmanager2/script/soundmanager2-nodebug-jsmin.js"></script>
        <!-- DateTimePicker -->
        <script src="/static/assets/vendor/jquery-datetimepicker/build/jquery.datetimepicker.full.min.js"></script>
        <!-- Revolution Slider -->
        <script src="/static/assets/vendor/revolution/js/jquery.themepunch.tools.min.js"></script>
        <script src="/static/assets/vendor/revolution/js/jquery.themepunch.revolution.min.js"></script>
        <script src="/static/assets/vendor/revolution/js/extensions/revolution.extension.video.min.js"></script>
        <script src="/static/assets/vendor/revolution/js/extensions/revolution.extension.carousel.min.js"></script>
        <script src="/static/assets/vendor/revolution/js/extensions/revolution.extension.navigation.min.js"></script>
        <!-- Keymaster -->
        <script src="/static/assets/vendor/keymaster/keymaster.js"></script>
        <!-- Summernote -->
        <script src="/static/assets/vendor/summernote/dist/summernote-bs4.min.js"></script>
        <!-- GODLIKE -->
        <script src="/static/assets/js/tibiakray.min.js"></script>
        <script src="/static/assets/js/tibiakray-init.js"></script>
		`
}
