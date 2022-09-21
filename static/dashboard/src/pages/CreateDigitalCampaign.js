import React, { Component } from "react";
import { Link, Redirect } from "react-router-dom"

import RichEditor from "../components/RichEditor"

import v1 from "../api/v1";
import InlineUpload from "../components/InlineUploader";
import PriceButtonGroup from "../components/payments/PriceButtonGroup";
import OptionsButtonGroup from "../components/OptionsButtongroup";

class CreateDigitalCampaign extends Component {
    constructor(props) {
        super(props)

        let tpe = this.props.template ? this.props.template : 'other'
        let role = this.props.type
        let allowedFileTypes = ['*/*']
        let maxFileSize = (1024* 1024 * 500)
        let maxFiles = 1

        switch (tpe){
            case 'video':
                allowedFileTypes  = ['video/*']
                maxFileSize = (Math.pow(1024, 3) * 2)
                maxFiles = 1
            break;
            case 'audio':
                allowedFileTypes = ['audio/mp3', 'audio/aac', 'audio/flac', 'audio/wav', 'audio/mpeg']
                maxFileSize = Math.pow(1024, 2) * 50
                maxFiles  = 1;
            break;
            case 'image':
                allowedFileTypes = ['image/*']
                maxFileSize = Math.pow(1024, 2) * 50
                maxFiles = 50
            break;
            default:
            break;
        }

        this.state = {
            published: '',
            multipleFiles: maxFiles > 1,
            maxNumberOfFiles: maxFiles,
            modalOpen: false,
            allowedFileTypes: allowedFileTypes,
            campaignimgurl: "/assets/img/upload.jpg",
            campaign: {
                type: tpe,
                title: '',
                price: 1.00,
                description: '  ',
                content: [],
                subscription: 'public',
                expires: 'never',
                download: true,
            }
        }

        this.onUploaded = (file) => {
            if (file._id) {
                let c = this.state.campaign;
                if(this.state.multipleFiles){
                    c.content.push(file._id)
                }else {
                    c.content = [file._id]
                }

                this.setState({
                    campaign: c,
                })
                if (! this.state.multipleFiles ){
                    this.setState({ campaignimgurl: v1.assets.imageURL(file._id, 480, 480)})
                }
            }
        }


        this.createNewCampaign = this.createNewCampaign.bind(this)
    }

    createNewCampaign() {
       
                if(this.state.campaign.content.length === 0){
                    return alert('You must upload at least one file.')
                }
        
                if(this.state.campaign.title.length < 3 ){
                    return alert('Title should be at least 3 characters long')
                }
        
                if (this.state.campaign.price < 0.50) {
                    return alert('Minimum service price is USD$0.50')
                }
        
                if(this.state.campaign.subscription === "pay_per_view"  ) {
                    let msg = ''
                    if ( this.state.campaign.price < 0.50 ){
                        msg = 'Minimum price for any content is USD$0.50'
                    }
                    if (this.state.campaign.price > 50.00) {
                        msg = 'Maximum price for any content is USD$50. Contact us for custom pricing.'
                    }
                    if(msg.length != 0){
                        return alert(msg)
                    }
                }
        
        v1.campaign.createCampaign(this.state.campaign, this.props.type).then(service => {
            console.log(service)
            if (service._id.length > 0 ) {
                //notify of creation success
                this.setState({published: '/creator/content/' + service._id})
            } else {
                alert('Failed to publish your content. Reason: '+ service.error)
            }
        }).catch(err => {
            alert('Failed to publish your content. Reason: \n'+ JSON.stringify(err))
        })
    }

    render() {
        let viewOptions = [
            {
                value: 'public',
                component: <><i class="fa fa-eye"></i> Public</> 
            },
            {
                value: 'pay_per_view',
                component: <><i class="fa fa-credit-card"></i> Pay per View</> 
            },
        ]

        if (this.props.user.subscriptions.active) {
            viewOptions.push({
                value: 'fans',
                component: <><i class="fa fa-lock"></i> Subscribers only</> 
            })
        }
        return this.state.published.length > 0 ? <Redirect to={this.state.published} /> : (

            <div className="row justify-content-center">
                <div class="col-12">
                    <h4>Upload your {this.props.type}</h4>
                </div>
                
                <div className="col-md-6" >
                    <div className="form-group">
                        <label>Content Title</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
                    <div className="form-group">
                        <label>Upload file(s)</label>
                        <InlineUpload type={this.props.template} maxNumberOfFiles={this.state.maxNumberOfFiles} purpose="content" onUploaded={this.onUploaded} allowedTypes={this.state.allowedFileTypes} />
                    </div>
                    
                    <div className="form-group">
                        <label>Description</label>
                        <RichEditor text={this.state.campaign.title} onChange={(text) => { let c = this.state.campaign; c.description = text; this.setState({ campaign: c }) }} />
                    </div>

                    <div className="form-group">
                        <label>Who can view this ?</label>
                        <OptionsButtonGroup items={viewOptions} item={this.state.campaign.subscription} onChange={(s) => { let c = this.state.campaign; c.subscription = s; this.setState({ campaign: c }) }} />
                    </div>
                    { this.state.campaign.subscription === "public" ? <></> :
                    <div className="form-group">
                            <label>Price per {this.state.campaign.subscription === 'fans' ? 'subscriber' : 'viewer'}</label>
                        <PriceButtonGroup prices={[0.5, 1, 2, 5, 10]} price={this.state.campaign.price} onChange={(price) => { let c = this.state.campaign; c.price = price; this.setState({ campaign: c }) }} />
                        
                    </div>
                    }

                    
                    <button className="btn btn-block btn-primary" onClick={this.createNewCampaign} ><i className="fa fa-check"></i> Create campaign</button>
                </div>
            </div>
        )
    }
}

export default CreateDigitalCampaign;