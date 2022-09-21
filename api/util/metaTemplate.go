package util

var metaTemplate string = `
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    <title>{{page_title}}</title>
    <!-- favicon -->

    
    <!-- Theme Color for Chrome, Firefox OS and Opera -->
    <meta name="theme-color" content="#4285f4">
	<meta name="generated_at" content="{{generated_at}}">
    
    <!-- Short description of the document (limit to 150 characters) -->
    <!-- This content *may* be used as a part of search engine results. -->
    <meta name="description" content="{{page_description}}">
    
    <!-- Control the behavior of search engine crawling and indexing -->
    <meta name="robots" content="index,follow"><!-- All Search Engines -->
    <meta name="googlebot" content="index,follow"><!-- Google Specific -->
    
    <!-- Tells Google not to show the sitelinks search box -->
    <meta name="google" content="nositelinkssearchbox">
    
    <!-- Tells Google not to provide a translation for this document -->
    <meta name="google" content="notranslate">
    
    <!-- Verify website ownership -->
    <meta name="google-site-verification" content="qoiVshYUbPqldqHWuYQ8Qm6s7VNI1Y2BrR0mN_v224w"><!-- Google Search Console -->
    <meta name="norton-safeweb-site-verification" content="norton_code"><!-- Norton Safe Web -->
    
    <!-- Identify the software used to build the document (i.e. - WordPress, Dreamweaver) -->
    <meta name="generator" content="myhustlezw">
        
    <!-- Gives a general age rating based on the document's content -->
    <meta name="rating" content="General">

    <meta property="fb:app_id" content="123456789">
    <meta property="og:url" content="{{page_url}}">
    <meta property="og:type" content="website">
    <meta property="og:title" content="{{page_title}}">
    <meta property="og:image" content="{{page_image}}">
    <meta property="og:image:alt" content="{{page_image_alt}}">
    <meta property="og:description" content="{{page_description}}">
    <meta property="og:site_name" content="MyHustle ZW">
    <meta property="og:locale" content="en_ZW">
    <meta property="article:author" content="{{page_author_url}}">

    <meta name="twitter:text:title" content="{{page_title}}"/>
    <meta name="twitter:image" content="{{page_image}}"/>
    <meta name="twitter:card" content="summary_large_image"/>
    <meta name="twitter:site" content="@myhustle_zw">
    <meta name="twitter:creator" content="{{twitter_username}}">
    <meta name="twitter:url" content="{{page_url}}">
    <meta name="twitter:title" content="{{page_title}}">
    <meta name="twitter:description" content="{{page_description}}">
    <meta name="twitter:image" content="{{page_image}}">

    <meta name="pinterest" content="nopin" description="Sorry, you can't save from my website!">
    

    <link rel=icon href=/favicon.ico sizes="20x20" type="image/png">

    
`
