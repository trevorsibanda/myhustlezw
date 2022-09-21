import React, { Component } from "react"
import { Link } from "react-router-dom"
import v1 from "../api/v1"

import ActionButtonUpload from "./ActionButtonUpload"

class PageCoverEditor extends Component {
    constructor(props) {
        super(props)

        this.state = {
            coverimageurl: v1.assets.imageURL(this.props.user.profile.cover_image, 0, 0),
            profpicurl: this.props.user.profile.profile_url
        }


        this.setCoverImage = (image) => {
            this.setState({
                coverimageurl: v1.assets.imageURL(image._id, 0, 0) + '?reload'
            })
        }
    }

    render() {
        return (
            <div className="row">
                <div className="col-xs-12 col-sm-12">
                    <div className="box box-inverse bg-img" style={{
                        "background-image": "url('" + this.state.coverimageurl + "')"
                    }} >
                        <div className="flexbox px-20 pt-20">
                            <label className="text-red">
                                <ActionButtonUpload noDashboard image={this.state.coverimageurl} type="image" onUploaded={this.setCoverImage} purpose="cover" allowedTypes={['image/*']} hidePreview={true} uploadBtnText="Upload cover" uploadBtnClasses="btn btn-info btn-sm " />
                                
                            </label>
                        </div>
                        <div className="box-body text-center pb-50">
                            <ActionButtonUpload noDashboard image={this.state.profpicurl} type="image" purpose="profile_pic" allowedTypes={['image/*']} hideButton={true} previewClasses="avatar avatar-xxl avatar-bordered avatar-upload-15vh" />
                            
                            <h4 className="mt-2 mb-0">
                                <a className="hover-primary text-white" href={"/@" + this.props.user.username} target="_blank">
                                    {this.props.user.fullname}
                                </a>
                            </h4>
                            <span style={{color: "black", backgroundColor: "white"}}>is {this.props.user.profile.description} <Link to="/creator/settings/mypage"><i class="fa fa-edit"></i></Link></span>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

export default PageCoverEditor;