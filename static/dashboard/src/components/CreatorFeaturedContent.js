import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'

import ReactPlayer from 'react-player'

import './common.css'



function CreatorFeaturedContent(props) {
    var config = {}
    if(props.featured.processed && !window.supportsHLS()) {
        config = { file: { forceHLS: true } }
    }
    return (
<section class="blog-area " id="player">
    <div class="single-blog-grid-01 ">
                <div class="thumb">
                    {props.featured.type === 'video' ?
                        <div className="player-wrapper">
                            <ReactPlayer config={config} controls={true} muted={false} width='100%' height='100%' url={props.stream_url} className="react-player" />
                        </div> : <></>}
                    {props.featured.type === 'image' ? 
                        <img src={props.featured.url} alt="featured" class="img-responsive" /> : <></>}
        </div>
    </div>
</section>
    )
}

export default CreatorFeaturedContent;