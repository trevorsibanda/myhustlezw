import React, { Component } from "react"
import { Link } from "react-router-dom"
import v1 from "../api/v1"
import RichEditor from "../components/RichEditor"

import ActionButtonUpload from "../components/ActionButtonUpload"
import PageCoverEditor from "../components/PageCoverEditor"
import Restrict from "../components/Restrict"
import PageDesignTips from "../components/PageDesignTips"

class BasicPageDetails extends Component {
    constructor(props) {
        super(props)

        this.state = {
            user: this.props.user,
            personal_website: this.props.user.profile.personal_website,
            currentFeatured: this.props.user.page.featured_url,
            currentProfilePic: this.props.user.profile.profile_url,
            canChangeUsername: false,
        }
        this.saveChanges = this.saveChanges.bind(this)

        this.setFeaturedImage = (image) => {
            this.setState({ currentFeatured: image.url })
        }

        this.setProfilePhoto = (image) => {
            this.setState({ currentProfilePic: image.url })
            
        }

        this.toggleCanEditUsername = () => {
             this.setState({canChangeUsername: ! this.state.canChangeUsername})
        }

    }

    saveChanges() {
        if (this.state.user.profile.description.length > 150) {
            return alert("'What are you creating?' field should be less than 150 characters.")
        }
        v1.user.updateBasicPageDetails({
            username: this.state.user.username,
            fullname: this.state.user.fullname,
            aboutme: this.state.user.profile.about_me,
            description: this.state.user.profile.description,
            ...this.state.personal_website,
        }).then(resp => {
            if (resp._id) {
                this.setState({ user: resp })
                alert("Successfully updated your page configuration.")
            } else {
                alert(resp.error)
            }
        })
    }


    render() {
        return (<>

            <h4 className="mt-0 mb-20">Basic details</h4>
            <Restrict user={this.props.user} >
                <div className="form-group">
                    <label>Profile photo</label>
                    <ActionButtonUpload image={this.state.currentProfilePic} allowedTypes={["image/*"]} onUploaded={this.setProfilePhoto} uploadBtnText="Upload profile photo" purpose="profile_pic" type="image" />
                </div>
            </Restrict>
            
            <div className="form-group">
                <label for="pwd">Username</label>
                <div className="input-group " data-children-count="1">
                    <input disabled={!this.state.canChangeUsername} readOnly={!this.state.canChangeUsername} type="text" className="form-control rounded" onChange={(evt) => { let u = this.state.user; u.username = evt.target.value; this.setState({ user: u }) }} value={this.state.user.username} />
                </div>
                <b style={{color: 'red'}}>All links to your content will break when you change your username. Please keep this in mind.</b>
                { !this.state.canChangeUsername ? 
                    <button className="btn btn-block btn-danger btn-sm" onClick={this.toggleCanEditUsername} ><i className="fa fa-link-external"></i> Change my username</button>
                    : <></>
                }
            </div>
            <div className="form-group">
                <label>Display name</label>
                <input type="text" className="form-control" placeholder="Enter full name" onChange={(evt) => { let u = this.state.user; u.fullname = evt.target.value; this.setState({ user: u }) }} value={this.state.user.fullname} />
            </div>
            <div className="form-group">
                <label>What are you creating ?</label>
                <p><small>This will be shown on the top of your page</small></p>
                <div class="input-group mb-3">
                    <textarea maxLength={151} rows={5} onChange={(evt) => { let user = this.state.user; user.profile.description = evt.target.value.substr(0, 175); this.setState({ user }) }}  class="form-control" placeholder="a content creator" >
                    {this.state.user.profile.description}
                    </textarea> 
                </div>
                <label><small> {this.state.user.fullname} is {this.state.user.profile.description} <span class={(this.state.user.profile.description.length <= 150 ? "text-success" : "text-bold text-danger")}>{this.state.user.profile.description.length}/150</span></small></label>
            </div>
            <div className="form-group">
                <label>Featured image/video</label>
                <ActionButtonUpload image={this.state.currentFeatured} allowedTypes={["image/*", "video/*"]} onUploaded={this.setFeaturedImage} uploadBtnText="Upload feature image/video" purpose="feature" type="image_video" />
            </div>
            <label>Personal website</label>
            <div class="input-group mb-3">
                <div class="input-group-prepend">
                    <span class="input-group-text"><i class="fa fa-link"></i></span>
                </div>
                <input type="url" maxLength={1024} class="form-control" placeholder="Leave blank to not show on your page" value={this.state.personal_website.url} onChange={evt => { let pw = this.state.personal_website; pw.url = evt.target.value; this.setState({personal_website: pw}) }} />
            </div>
            <label>Facebook</label>
            <div class="input-group mb-3">
                <div class="input-group-prepend">
                    <span class="input-group-text"><i class="fa fa-facebook"></i></span>
                </div>
                <input type="url" maxLength={200} class="form-control" placeholder="Leave blank to not show on your page" value={this.state.personal_website.facebook} onChange={evt => { let pw = this.state.personal_website; pw.facebook = evt.target.value; this.setState({personal_website: pw}) }} />
            </div>
            <label>Twitter Username</label>
            <div class="input-group mb-3">
                <div class="input-group-prepend">
                    <span class="input-group-text"><i class="fa fa-twitter"></i></span>
                </div>
                <input type="text" maxLength={15} class="form-control" placeholder="Leave blank to not show on your page" value={this.state.personal_website.twitter} onChange={evt => { let pw = this.state.personal_website; pw.twitter = evt.target.value; this.setState({personal_website: pw}) }} />
            </div>
            <label>Instagram Username</label>
            <div class="input-group mb-3">
                <div class="input-group-prepend">
                    <span class="input-group-text"><i class="fa fa-instagram"></i></span>
                </div>
                <input type="url" maxLength={ 30 } class="form-control" placeholder="Leave blank to not show on your page" value={this.state.personal_website.instagram} onChange={evt => { let pw = this.state.personal_website; pw.instagram = evt.target.value; this.setState({personal_website: pw}) }} />
            </div>
            <label>Youtube </label>
            <div class="input-group mb-3">
                <div class="input-group-prepend">
                    <span class="input-group-text"><i class="fa fa-youtube text-google"></i></span>
                </div>
                <input type="url" maxLength={100} class="form-control" placeholder="Leave blank to not show on your page" value={this.state.personal_website.youtube} onChange={evt => { let pw = this.state.personal_website; pw.youtube = evt.target.value; this.setState({personal_website: pw}) }} />
            </div>
            <div className="row justify-content-center" >
                <div className="col-md-8" >
                    <button className="btn btn-block btn-primary" onClick={this.saveChanges} ><i className="fa fa-check"></i> Save changes</button>
                </div>
            </div>

        </>)
    }
}


class PageSettings extends Component {
    render() {
        return (
            <div className="row">
                <div className="col-lg-6 col-md-7">
                    <BasicPageDetails user={this.props.user} />
                </div>
                <div className="col-lg-6 col-md-5 " >
                    <PageDesignTips user={this.props.user} />
                </div>
            </div>
        )
    }
}

export default PageSettings;