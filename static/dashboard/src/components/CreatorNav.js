import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'
import Linkify from 'react-linkify'

class CreatorTopNav extends Component {
    constructor(props) {
        super(props)

        this.closeMenu = () => {
            document.getElementById('mobileNavBtn').click()
        }
    }

    render() {
        let sm = this.props.creator.profile.personal_website
        return (
        <div class="box margin-top-70 bg-img">
            <div class="flexbox align-items-center px-20" data-overlay="4" style={{backgroundColor: "ghostwhite", color: "black"}}>
                <div class="flexbox align-items-center mr-auto" >
                <Link to={"/@" + this.props.creator.username} >
                            <img class="avatar avatar-xl avatar-bordered" style={{ minWidth: "64px" }} src={this.props.creator.profile.profile_url} alt="" />
                </Link>
                <div class="pl-10 padding-top-20">
                <h4>
                    <Link to={"/@" + this.props.creator.username} class="hover-primary text-black" href="#">{ this.props.creator.fullname }</Link>
                </h4>
                <span class="">
                    <Linkify
                        componentDecorator={(decoratedHref, decoratedText, key) => (
                            <a target="blank" style={{color: 'red', fontWeight: 'bold'}} rel="noopener" target="_blank" href={decoratedHref} key={key}>
                                {decoratedText}
                            </a>
                        )}
                    >{ this.props.creator.profile.description }</Linkify>
                </span>
                <p>
                </p>
                    <div class="gap-items font-size-16">
                    {sm.url && sm.url.length > 1 ? <a class="text-default" href={sm.url} target="_blank" rel="noopener"><i class="fa fa-url"></i></a> : <></>}
                    {sm.facebook && sm.facebook.length > 1 ? <a class="text-facebook" href={sm.facebook} target="_blank" rel="noopener"><i class="fa fa-facebook"></i></a> : <></>}
                    {sm.instagram && sm.instagram.length > 1 ? <a class="text-instagram" href={"https://instagram.com/" + sm.instagram} target="_blank" rel="noopener"><i class="fa fa-instagram"></i></a> : <></>}
                    {sm.google && sm.google.length > 1 ? <a class="text-google" href={sm.youtube} target="_blank" rel="noopener"><i class="fa fa-youtube"></i></a> : <></>}
                    {sm.twitter && sm.twitter.length > 1 ? <a class="text-twitter" href={"https://twitter.com/"+ sm.twitter} rel="noopener"><i class="fa fa-twitter"></i></a> : <></>}
                </div>
                <p></p>
            </div>

        </div>
    </div>

</div>
        );
    }
}

export default CreatorTopNav;