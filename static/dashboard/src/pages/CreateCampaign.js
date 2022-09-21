import { ContentBlock } from "draft-js";
import React, { Component } from "react";
import { Link, NavLink } from "react-router-dom"
import v1 from "../api/v1";

import Preloader from "../components/PreLoader";
import Restrict from "../components/Restrict";


import CreateDigitalCampaign from "../pages/CreateDigitalCampaign"
import CreateServiceCampaign from "../pages/CreateServiceCampaign"
import CreateYoutubeEmbedCampaign from "./CreateYoutubeEmbedCampaign";


class CreateNewCampaignPicker extends Component {
    constructor(props){
        super(props)

        this.state = {
            service_templates: []
        }

        v1.config.serviceTemplates().then(templates => {
            if(templates.length > 0){
                this.setState({ service_templates: templates })
            }
        })
    }
    render(){
        return (
            <>
            <h4>Upload your content</h4>
            <div class="row">
                
                
                <div class="col-lg-3 col-6">
                    <a class="box box-link-shadow text-center" href="#" onClick={() => this.props.setType('content', 'image')}>
                        <div class="box-body">
                            <div class="font-size-24">Photos</div>
                            <span>Upload one or more images into one collection.</span>
                        </div>
                        <div class="box-body bg-warning text-white">
                            <p>
                                <span class="fa fa-camera font-size-30 text-white"></span>
                            </p>
                        </div>
                    </a>
                </div>
                <div class="col-lg-3 col-6">
                    <a class="box box-link-shadow text-center" href="#" onClick={()=> this.props.setType('content', 'video')}>
                        <div class="box-body">
                            <div class="font-size-24">Video</div>
                            <span>Upload a video which can be streamed from your page.</span>
                        </div>
                        <div class="box-body bg-info text-white">
                            <p>
                                <span class="fa fa-film font-size-30 text-white"></span>
                            </p>
                        </div>
                    </a>
                </div>
                <div class="col-lg-3 col-6">
                    <a class="box box-link-shadow text-center" href="#" onClick={() => this.props.setType('embed', 'youtube')}>
                        <div class="box-body">
                            <div class="font-size-24">Youtube video</div>
                            <span>Embed a video from Youtube here.</span>
                        </div>
                        <a class="box-body bg-google text-white">
                            <span class="fa fa-youtube font-size-30" ></span>
                        </a>
                    </a>
                </div>
                <div class="col-lg-3 col-6">
                        <a class="box box-link-shadow text-center" href="#" onClick={() => { this.props.setType('service', {})}}>
                        <div class="box-body">
                            <div class="font-size-24">Offer a service</div>
                            <span>Sell subscriptions, private access, consultation... etc</span>
                        </div>
                        <div class="box-body bg-success ">
                            <p>
                                <span class="fa fa-archive text-white font-size-30"></span>
                            </p>
                        </div>
                    </a>
                </div>
            </div>
                 
            </>
        )
    }
}


class CreateNewCampaign extends Component {
    constructor(props) {
        super(props)
        this.state = {
            type: 'choose',
            template: {},
            templates: [],
            loading: false,
        }


        this.pickType = this.pickType.bind(this)

    }

    pickType(tpe, template) {
        this.setState({type: tpe, template})
        console.log(this.state)
    }

    render() {
        let content = <></>
        switch (this.state.type) {
            case 'choose':
                content = <CreateNewCampaignPicker setType={this.pickType} />
                break;
            case 'content':
                content = <Restrict user={this.props.user} showWarning={true}>
                    <CreateDigitalCampaign user={this.props.user} type={this.state.type} template={this.state.template} />
                    </Restrict>
            break;
            case 'embed':
                content = <Restrict user={this.props.user} showWarning={true}><CreateYoutubeEmbedCampaign user={this.props.user} /></Restrict>
            break;
            case 'service':
                content = <Restrict user={this.props.user} showWarning={true} >
                    <CreateServiceCampaign user={this.props.user} template={this.state.template} />
                    </Restrict>
                break;
            default:
                content = <p>Internal error: Unknown internal state</p>
        }

        if (this.state.type !== 'choose' ){
            content = <>
                <div class="content-header">
                    <div class="mr-auto">
                    <div class="d-inline-block align-items-center">
                        <nav>
                            <ul class="breadcrumb fa-2x">
                                <li class="breadcrumb-item"><a href="javascript:;" onClick={_ => this.setState({ type: 'choose' })}><i class="fa fa-plus"></i> Create new</a>
                                /{this.state.type}
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>
                </div>
                {content}</>
        }

        return this.state.loading ? <Preloader /> : content;
    }

}

export default CreateNewCampaign;