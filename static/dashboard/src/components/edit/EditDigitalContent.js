import React, { Component } from "react";

import RichEditor from "../RichEditor"

import v1 from "../../api/v1";
import PriceButtonGroup from "../payments/PriceButtonGroup";
import OptionsButtonGroup from "../OptionsButtongroup";
import InlineUpload from "../InlineUploader";
import LiteYouTubeEmbed from "react-lite-youtube-embed";
import ActionButtonUpload from "../ActionButtonUpload";
import ReactPlayer from "react-player";

class EditDigitalContent extends Component {
    constructor(props) {
        super(props)

        v1.page.event('Edit Content', 'Edit', props.content.type)

        let viewOptions = [
            {
                value: 'public',
                component: <><i class="fa fa-eye"></i> Public</> 
            }]
        
        if (props.content.type !== 'embed' && props.content.type !== 'service') {
            viewOptions.push({
                value: 'pay_per_view',
                    component: <><i class="fa fa-credit-card"></i> PayPerView</>
            })
        }
        
        if (this.props.creator.subscriptions.active && props.content.type !== 'embed') {
            viewOptions.push({
                value: 'fans',
                component: <><i class="fa fa-lock"></i> Subscribers</> 
            })
        }
        
        this.state = {
            campaign: this.props.content,
            youtube_id: '',
            youtube_url: 'https://youtube.com/watch?v=' + (props.content && props.content.remote_id ? props.content.remote_id : '' ),
            valid_yt_link: false,
            viewOptions,
        }

        this.updateCampaign = this.updateCampaign.bind(this)
        this.validateEmbed = this.validateEmbed.bind(this)
        this.validatePayPerView = this.validatePayPerView.bind(this)
        this.validateUpload = this.validateUpload.bind(this)
        this.validateService = this.validateService.bind(this)

        this.serviceForm = this.serviceForm.bind(this)
        this.embedForm = this.embedForm.bind(this)
        this.generalContentForm = this.generalContentForm.bind(this)


        this.getYoutubeID = (link) => {
            return link.match(/(?:youtu\.be\/|youtube\.com(?:\/embed\/|\/v\/|\/watch\?v=|\/user\/\S+|\/ytscreeningroom\?v=|\/sandalsResorts#\w\/\w\/.*\/))([^\/&]{10,12})/)[1];
        }

        this.applyYoutubeLink = (link) => {
            try{
                let id = this.getYoutubeID(link)
                if (id.length > 0) {
                    this.setState({youtube_id: id, valid_yt_link: true, youtube_url: link})
                    let c = this.state.campaign
                    c.content = [id]
                    c.remote_id = id
                    this.setState({campaign: c})
                }
            }catch(e){
                this.setState({valid_yt_link: false, youtube_url: link})
                console.log(e)
            }
        }

        this.onUploaded = (file) => {
            if (file._id) {
                let c = this.state.campaign;
                c.preview = file._id
                c.preview_url = file.url

                this.setState({campaign: c })
            }
        }

        if (props.content.type === 'embed') {
            this.applyYoutubeLink(this.state.youtube_url)
        }
    }

    validateEmbed() {
        this.applyYoutubeLink(this.state.youtube_url)
        if (!this.state.valid_yt_link) {
            alert('Please enter a valid youtube link')
            return false
        }
        if (this.state.campaign.title.length <= 3) {
            alert('Title must be at least 3 characters long')
            return false
        }
        return true
    }

    validatePayPerView() {
        var c = this.state.campaign

        if (c.title.length < 10) {
            alert('Title must be at least 10 characters long')
            return false
        }

        if (c.description.length < 10) {
            alert('Description must be at least 10 characters long')
            return false
        }

        if (c.price < 0.5 || c.price > 50) {
            alert('Pay per view prices can only be between $0.50 and $10')
            return false
        }

        return true
    }

    validateService() {
        var c = this.state.campaign
        if (c.title.length < 10) {
            alert('Title must be at least 10 characters long')
            return false
        }

        if (c.description.length < 10) {
            alert('Description must be at least 10 characters long')
            return false
        }

        if (c.price < 1 || c.price > 50) {
            alert('Service prices can only be between $1 and $50')
            return false
        }
       

        if(c.service.quantity_available > 100) {
            alert('Maximum quantity is 100. You can always increase this number back to 100 when it gets low')
            return false
        }

        if (c.service.question.length < 5) {
            alert('To collect information from the paer you must ask a question. Min length is 5 characters')
            return false
        }

        return true

    }

    validateUpload() {
        var c = this.state.campaign

        if (c.title.length < 10) {
            alert('Title must be at least 10 characters long')
            return false
        }

        if (c.description.length < 10) {
            alert('Description must be at least 10 characters long')
            return false
        }

        return true

    }


    serviceForm() {
        return (
            <>
            <div className="form-group">
                <label>What are you offering?</label>
                <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => {let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c })}} />
                
            </div>
            
            <div className="form-group">
                <label>Image to show</label>
                <ActionButtonUpload image={this.state.campaign.preview_url} allowedTypes={["image/*"]} onUploaded={this.onUploaded} uploadBtnText="Upload new service image" purpose="service_preview" type="image" />
                <p>
                    <small>Image will be shown on the service page. This can be a poster, an example or instructions.</small>
                </p>
            </div>
            <div className="form-group">
                <label>Description/Instructions</label>
                <p>
                    <small>Provide details or instructions to your users before they pay for your service. 
                        You can highlight your terms and conditions here.</small>
                </p>
                <RichEditor text={this.state.campaign.service.instructions} onChange={(text) => { let c = this.state.campaign; c.service.instructions = text; this.setState({ campaign: c }) }} />
            </div>
            <h5>Collect details</h5>
            <p>Ask a question before a customer can order your service.</p>
            <div className="form-group">
                <label>Ask a question</label>
                <input type="text" class="form-control" placeholder="Ask the payer a question" value={this.state.campaign.service.question} onChange={(evt) => { let c = this.state.campaign; c.service.question= evt.target.value; this.setState({ campaign: c }) }}/>
            </div>
            <br/>
            <div class="form-group" >
                <label>Price</label>
                <PriceButtonGroup prices={[1,2,5,10,20,50]} price={this.state.campaign.price} onChange={(price) => {let c = this.state.campaign; c.price = price; this.setState({ campaign: c })}} />
                <p><small>All prices are listed in USD. Payments in ZWL will be crossrated at the official rate.</small></p>
            </div>
            <div className="form-group">
                <label>Thank you message</label>
                <RichEditor text={this.state.campaign.service.thankyou_message} onChange={(text) => { let c = this.state.campaign; c.service.thankyou_message = text; this.setState({ campaign: c }) }} />
            </div>
            
            
            <h5>Restrictions</h5>
            <div class="row">
                <div class="col-md-6">
                    <div className="form-group">
                        <label>Limit available slots</label>
                        <input type="number" min={0} max={100} placeholder="Maximum available orderd you can take" class="form-control" value={this.state.campaign.service.quantity_available} onChange={(evt) => { let c = this.state.campaign; c.service.quantity_available = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
                </div>
            </div>
            
            
            </>
        )
    }

    embedForm() {
        return(
            <div className="row">
                <div class="col-12">
                    <h4>Embed a youtube video.</h4>
                    <p><small>Already have content on Youtube? Just embed it!</small></p>
                </div>
                <div className="col-md-12" >
                    <div className="form-group">
                        <label>Video Title</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
                    
                </div>
                <div className="col-md-12" >
                    <div className="player-wrapper">
                        <ReactPlayer width='100%' controls={true} height='100%' playing={true} url={"https://youtu.be/" + this.state.youtube_id} className="react-player" />
                    </div>
                    <p>
                        <div className="form-group">
                            <label>Youtube Video Link</label>
                            <div className="input-group mb-3">
                                <div className="input-group-prepend">
                                    <span className="input-group-text"><i class="fa fa-youtube text-google"></i></span>
                                </div>
                                <input type="text" className="form-control" placeholder="https://www.youtube.com/watch?v=BkSdD5VtyRM" value={this.state.youtube_url} onChange={(evt) => this.applyYoutubeLink(evt.target.value)} />
                            </div>
                        </div>
                        <small>WARNING: Only include videos which you own the legal copyright to if you plan to monetize from the content. 
                            Including a video you do not legally own and monetizing it will result in your account being suspended.</small>
                    </p>
                </div>
                <div className="col-md-12" >
                    
                    <div className="form-group">
                        <label>Who can view this ?</label>
                        <OptionsButtonGroup items={this.state.viewOptions} item={this.state.campaign.subscription} onChange={(s) => { let c = this.state.campaign; c.subscription = s; this.setState({ campaign: c }) }} />
                        <p><small>Embedded content can only be public.</small></p>
                    </div>
                </div>
            </div>
        )
    }

    generalContentForm() {
        return (
            <div className="row justify-content-center">
                <div className="col-md-12" >
                    <div className="form-group">
                        <label>Content Title</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
                    {this.props.createMode ?
                        <div className="form-group">
                            <label>Upload file(s)</label>
                            <InlineUpload type={this.props.template} maxNumberOfFiles={this.state.maxNumberOfFiles} purpose="content" onUploaded={this.onUploaded} allowedTypes={this.state.allowedFileTypes} />
                    </div> : <></>}
                    
                    <div className="form-group">
                        <label>Description</label>
                        <RichEditor text={this.state.campaign.title} onChange={(text) => { let c = this.state.campaign; c.description = text; this.setState({ campaign: c }) }} />
                    </div>

                    <div className="form-group">
                        <label>Who can view this ?</label>
                        <OptionsButtonGroup items={this.state.viewOptions} item={this.state.campaign.subscription} onChange={(s) => { let c = this.state.campaign; c.subscription = s; this.setState({ campaign: c }) }} />
                    </div>
                    { this.state.campaign.subscription === "public" ? <></> :
                    <div className="form-group">
                            <label>Price per {this.state.campaign.subscription === 'fans' ? 'subscriber' : 'viewer'}</label>
                        <PriceButtonGroup prices={[0.5, 1, 2, 5, 10]} price={this.state.campaign.price} onChange={(price) => { let c = this.state.campaign; c.price = price; this.setState({ campaign: c }) }} />
                        
                    </div>
                    }

                </div>
            </div>
        )
    }

    updateCampaign() {

        var fn
        var c = this.state.campaign
        if (c.type === 'service') {
            fn = this.validateService
        } else if (c.subscription === 'pay_per_view') {
            fn = this.validatePayPerView
        } else if (c.subscription === 'public') {
            if (c.type === 'embed') {
                fn = this.validateEmbed
            } else {
                fn = this.validateUpload
            }
        }
        
        if (!fn()) {
            return
        }

        v1.campaign.update(this.state.campaign).then(resp => {
            if (resp && resp.campaign && resp.campaign._id) {
                this.setState({ campaign: resp.campaign })
                alert('Successfully saved the changes.')
                return(this.props.onSave ? this.props.onSave(resp.campaign) : true )
            } else {
                alert('Failed to save changes. Reason: '+ (resp.error ? resp.error : resp))
            }
        }).catch(err => {
            alert('Error: ' + err)
        })
    }

    render() {
        var component
        switch (this.state.campaign.type) {
            case 'service': component = this.serviceForm; break;
            case 'embed': component = this.embedForm; break;
            default: component = this.generalContentForm; break;
        }
        
        return (

            <div className="row justify-content-center">
                <div className="col-md-12 mt-20" >
                    {component()}
                    <button className="btn btn-block btn-primary" onClick={this.updateCampaign} ><i className="fa fa-check"></i> Save changes</button>
                </div>
            </div>
        )
    }
}

export default EditDigitalContent;