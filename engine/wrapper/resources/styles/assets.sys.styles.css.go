package styles

var File_assets_sys_styles_css = []byte(`html {
    min-height: 100%;
    height: 100%;
    position: relative
}

body {
    background: #205081 url(/assets/sys/bg.png) repeat 0 0;
    margin: 0;
    padding: 0;
    color: #FFF;
    font-family: Arial, sans-serif;
    font-weight: 300;
    font-size: 18px;
    line-height: 21px;
    position: relative;
    width: 100%;
    min-height: 100%;
    height: 100%;
    height: auto;
    display: table
}

.wrapper {
    padding: 15px;
    text-align: center;
    display: table-cell;
    vertical-align: middle
}

.wrapper .logo {
    width: 332px;
    height: 203px;
    margin: 0 auto;
    background: url(/assets/sys/logo.png) no-repeat 0 0
}

.wrapper .logo .svg {
    width: 150px;
    height: 150px;
    margin: 0 auto;
    position: relative;
    padding-top: 40px
}

.wrapper .logo .svg img {
    top: 50%;
    left: 50%;
    position: absolute;
    width: 150px;
    height: 150px;
    margin-top: -75px;
    margin-left: -75px;
    animation-name: fave;
    animation-duration: 1000ms;
    animation-iteration-count: infinite;
    animation-timing-function: linear
}

@keyframes fave {
    0% {
        width: 150px;
        height: 150px;
        margin-top: -75px;
        margin-left: -75px
    }

    40% {
        width: 150px;
        height: 150px;
        margin-top: -75px;
        margin-left: -75px
    }

    60% {
        width: 120px;
        height: 120px;
        margin-top: -60px;
        margin-left: -60px
    }

    100% {
        width: 150px;
        height: 150px;
        margin-top: -75px;
        margin-left: -75px
    }
}

h1, h2 {
    font-weight: 400;
    font-size: 34px;
    line-height: 36px;
    margin: 10px 0
}

h2 {
    font-weight: 350;
    font-size: 14px;
    line-height: 18px;
    margin-bottom: 0
}

@media only screen and (max-width:800px) {
    .wrapper {
        padding: 15px 0
    }

    h1, h2 {
        padding: 0 15px
    }
}`)
