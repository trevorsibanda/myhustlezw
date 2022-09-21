import React, { Component } from "react";
import { Link, useHistory } from "react-router-dom"

import RichEditor from "../components/RichEditor"

import v1 from "../api/v1";
import InlineUpload from "../components/InlineUploader";
import PriceButtonGroup from "../components/payments/PriceButtonGroup";

class CreateServiceCampaign extends Component {
    constructor(props) {
        super(props)
        this.createNewService = this.createNewService.bind(this)
        this.state = {
            campaignimgurl: "https://img2.goodfon.com/wallpaper/nbig/d/31/art-gorod-budushie-razvaleny.jpg",
            campaign: {
                type: 'service',
                title: this.props.template.title ? this.props.template.title : '',
                price: this.props.template.price ? this.props.template.price : 1.00,
                description: this.props.template.description ? this.props.template.description : '',
                preview: '',
                subscription: 'public',
                question: this.props.template.question ? this.props.template.question : '',
                thankyou: this.props.template.thankyou ? this.props.template.thankyou : 'Thank you message here',
                instructions: this.props.template.instructions ? this.props.template.instructions : '',
                quantity: this.props.template.quantity ? this.props.template.quantity : '',
            }
        }

        this.onUploaded = (file) => {
            if (file._id) {
                let c = this.state.campaign;
                c.preview = file._id

                this.setState({ campaignimgurl: v1.assets.imageURL(file._id, 480, 480), campaign: c })
            }
        }
    }

    createNewService() {
        //validate before creating

        if(this.state.campaign.preview === ''){
            return alert('No preview image upload. Upload one image please.')
        }

        if(this.state.campaign.title.length < 10 ){
            return alert('Title should be at least 10 characters long')
        }

        if (this.state.campaign.price < 1.00) {
            return alert('Minimum service price is USD$1.00')
        }

        if(this.state.campaign.instructions.length < 10) {
            return alert('Minimum instructions should be 10 characters long.')
        }

        if(this.state.campaign.thankyou.length < 50) {
            return alert('The thank you message should be at least 50 characters long.')
        }

        v1.campaign.createService(this.state.campaign).then(service => {
            if(service._id){
                //notify of creation success
                useHistory().push('/creator/content/'+ service._id)
            }else if(service.error) {
                alert('Failed to create new service. Error: ' + service.error)
            }
        })
    }

    render(){
        return (
            <div className="row justofy-content-center">
                <div class="col-12">
                    <h4>Offer a service</h4>
                </div>
                <div className="col-md-10" >
                    <div className="form-group">
                        <label>What are you offering?</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => {let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c })}} />
                        
                    </div>
                    
                    <div class="form-group" >
                        <label>Price</label>
                        <PriceButtonGroup prices={[1,2,5,10,20,50]} price={this.state.campaign.price} onChange={(price) => {let c = this.state.campaign; c.price = price; this.setState({ campaign: c })}} />
                        <p><small>All prices are listed in USD. Payments in ZWL will be crossrated at the official rate.</small></p>
                    </div>
                    <div className="form-group">
                        <label>Description/Instructions</label>
                        <p>
                            <small>Provide details or instructions to your users before they pay for your service. 
                                You can highlight your terms and conditions here.</small>
                        </p>
                        <RichEditor text={this.state.campaign.instructions} onChange={(text) => { let c = this.state.campaign; c.instructions = text; this.setState({ campaign: c }) }} />
                    </div>
                    <div className="form-group">
                        <label>Image to show</label>
                        <InlineUpload type="image" maxNumberOfFiles={1} purpose="service_preview" onUploaded={this.onUploaded} allowedTypes={['image/*']} />
                        <p>
                            <small>Image will be shown on the service page. This can be a poster, an example or instructions.</small>
                        </p>
                    </div>
                    <div className="form-group">
                        <label>Thank you message</label>
                        <RichEditor text={this.state.campaign.thankyou} onChange={(text) => { let c = this.state.campaign; c.thankyou = text; this.setState({ campaign: c }) }} />
                    </div>
                    <h5>Collect details</h5>
                    <p>Ask a question or offer options before a customer can order your service.</p>
                    <div className="form-group">
                        <label>Ask a question</label>
                        <input type="text" class="form-control" placeholder="Ask the user a question" value={this.state.campaign.question} onChange={(evt) => { let c = this.state.campaign; c.question= evt.target.value; this.setState({ campaign: c }) }}/>
                    </div>
                    <br/>
                    <h5>Restrictions</h5>
                    <div class="row">
                        <div class="col-md-6">
                            <div className="form-group">
                                <label>Limit available slots</label>
                                <input type="number" placeholder="Maximum available orderd you can take" class="form-control" value={this.state.campaign.quantity} onChange={(evt) => { let c = this.state.campaign; c.quantity = evt.target.value; this.setState({ campaign: c }) }} />
                            </div>
                        </div>
                    </div>
                    
                    <button className="btn btn-block btn-primary" onClick={this.createNewService} ><i className="fa fa-check"></i> Create service</button>
                </div>
            </div>

        )
    }
}

export default CreateServiceCampaign;