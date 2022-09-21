import React, { Component } from "react";
import { Link, Route, Switch } from "react-router-dom"
import v1 from "../api/v1";
import ContentPreviewCard from "../components/ContentPreviewCard";
import SupportersList from "../components/ListSupporters";
import RichEditor from "../components/RichEditor"

class CampaignSummary extends Component {
    render(){
        
        let files = this.props.files ? this.props.files : []
        let file =  files && files.length > 0 ? files[0] : null
        return (
             <>
                {files.length > 0 ?
            <ContentPreviewCard 
                file={file}
                title={this.props.campaign.title}
                description={this.props.campaign.description}
                link={'/@'+this.props.user.username + '/'+ this.props.campaign.uri}
                files={files}
                campaign={this.props.campaign}
                visibility={this.props.campaign.subscription}
            /> : <img src="/assets/blank.png" class="img-responsive" /> }
            
            </>
        )
    }
}

class CampaignRevenueSummary extends Component {
    render(){
        return (
            <>
            <h2>Revenue</h2>
            </>
        )
    }
}

class EditServiceCampaign extends Component{
    constructor(props){
        super(props)
        this.state = {
            campaign: this.props.campaign,
        }

        this.updateCampaign = this.updateCampaign.bind(this)
    }

    updateCampaign() {
        v1.campaign.update(this.state.campaign).then(resp => {
            if (resp.status === 'ok') {
                this.setState({ campaign: resp.campaign })
                alert("Successfully saved changes :)")
            } else {
                alert(resp.error)
            }
        }).catch(alert)
    }

    render(){
        return (
            <>
                <div className="form-group">
                    <label>What are you offering?</label>
                    <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                    <p>
                        <small>Your campaign will be available at <a href="" target="_blank">https://myhustle.co.zw/@trevorsibb/my-intro-message</a> </small>
                    </p>
                </div>

                <div class="form-group" >
                    <label>Price</label>
                    <div class="input-group" data-children-count="1">
                        <div class="input-group-addon">USD $</div>
                        <input type="number" class="form-control" value={this.state.campaign.price} onChange={(evt) => { let c = this.state.campaign; c.price = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
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
                    <label>Featured image/video</label>
                    <a href="#" onClick={this.handleOpen} >
                        <img src={this.state.campaignimgurl}
                            class="img-thumb"
                        />
                    </a>
                    <p>
                        <small>Featured image or video will be shown on the account page.</small>
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
                    <input type="text" class="form-control" placeholder="Ask the user a question" value={this.state.campaign.question} onChange={(evt) => { let c = this.state.campaign; c.question = evt.target.value; this.setState({ campaign: c }) }} />
                </div>
                <br />
                <h5>Restrictions</h5>
                <div class="row">
                    <div class="col-md-6">
                        <div className="form-group">
                            <label>Limit available slots</label>
                            <input type="number" placeholder="Maximum available orderd you can take" class="form-control" value={this.state.campaign.quantity} onChange={(evt) => { let c = this.state.campaign; c.quantity = evt.target.value; this.setState({ campaign: c }) }} />
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div className="form-group">
                            <label>Expires after</label>
                            <select className="form-control" value={this.state.campaign.expires} onChange={(evt) => { let c = this.state.campaign; c.expires = evt.target.value; this.setState({ campaign: c }) }}>
                                <option value="never">Never</option>
                                <option value="1 year">1 year</option>
                                <option value="3 months">3 months</option>
                                <option value="1 month">1 month </option>
                                <option value="1 week">1 week</option>
                                <option value="3 days">3 days</option>
                                <option value="1 day">1 day</option>
                            </select>
                        </div>
                    </div>
                </div>
                <button className="btn btn-block btn-primary" onClick={this.updateCampaign} ><i className="fa fa-check"></i> Save changes</button>

            </>
        )
    }
}

class EditDigitalCampaign extends Component {
    constructor(props){
        super(props)
        this.state = {
            campaign: this.props.campaign,
        }

        this.updateCampaign = this.updateCampaign.bind(this)
    }

    updateCampaign(){
        v1.campaign.update(this.state.campaign).then(resp => {
            if(resp.status === 'ok'){
                this.setState({ campaign: resp.campaign })
            }else{
                alert(resp.error)
            }
        }).catch(alert)
    }

    render(){
        return (
            <>
                <div className="" >
                    <div className="form-group">
                        <label>Content Title</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                        <p>
                            <small>Your campaign will be available at <a href="" target="_blank">https://myhustle.co.zw/@trevorsibb/my-intro-message</a> </small>
                        </p>
                    </div>
                    <div className="form-group">
                        <label>Who can view this ?</label>
                        <select className="form-control" value={this.state.campaign.subscription} onChange={(evt) => { let c = this.state.campaign; c.subscription = evt.target.value; this.setState({ campaign: c }) }}>
                            <option value="public">Everyone (Public )</option>
                            {this.props.user.payments_active ? <>
                                <option value="fans">My fans only ( Subscription ) </option>
                            </> : <></>}
                            {this.props.user.payments_active ? <option value="pay_per_view">Pay per view/download ( For Sale)</option> : <></>
                            }
                        </select>
                        {this.props.user.payments_active ? <></> : <div class="alert alert-warning">Verify your account to create paid content</div>}
                    </div>
                    {!this.props.user.payments_active || this.state.campaign.subscription != 'pay_per_view' ? <></> :
                        <div className="form-group">
                            <label>Price per item</label>
                            <div className="input-group mb-3">
                                <div className="input-group-prepend">
                                    <span className="input-group-text">USD $</span>
                                </div>
                                <input type="number" className="form-control" placeholder="0.00" value={this.state.campaign.price} onChange={(evt) => { let c = this.state.campaign; c.price = evt.target.value; this.setState({ campaign: c }) }} />
                            </div>
                        </div>
                    }
                    <div className="form-group">
                        <label>Description</label>
                        <RichEditor text={this.state.campaign.description} onChange={(text) => { let c = this.state.campaign; c.description = text; this.setState({ campaign: c }) }} />
                    </div>
                    <div className="form-group">
                        <label>Expires after</label>
                        <select className="form-control" value={this.state.campaign.expires} onChange={(evt) => { let c = this.state.campaign; c.expires = evt.target.value; this.setState({ campaign: c }) }}>
                            <option value="never">Never</option>
                            <option value="1 year">1 year</option>
                            <option value="3 months">3 months</option>
                            <option value="1 month">1 month </option>
                            <option value="1 week">1 week</option>
                            <option value="3 days">3 days</option>
                            <option value="1 day">1 day</option>
                        </select>
                    </div>
                    <button className="btn btn-block btn-primary" onClick={this.updateCampaign} ><i className="fa fa-check"></i> Save changes</button>
                </div>
            </>
        )
    }
}

class EditEmbedCampaign extends Component {
    render(){
        return <>Not implemented</>
    }
}

class CampaignOrders extends Component {
    render() {
        return <>Not implemented</>
    }
}

class CampaignSupporters extends Component {
    render() {
        return <>Not implemented</>
    }
}

class CampaignEngagement extends Component {
    render() {
        return <>Not implemented</>
    }
}


class CampaignNotifications extends Component {
    render() {
        return <>Not implemented</>
    }
}

class CampaignDelete extends Component {
    render() {
        return <>Not implemented</>
    }
}


class EditCampaignLandingPage extends Component {

    
}

class EditCampaign extends Component {

    constructor(props) {
        super(props)

        let { id } = this.props.match.params
        

        this.state = {
            id: id,
            campaign: {},
            files: [],
            file: {},
            supporters: [],
            stats: {},
            user: {},
            type: '',
            loading: true
        }


        

        this.pages = () => [
            { title: "Edit content", icon:"edit", lit: false, id: "edit"},
            { title: "Supporters", icon: "heart text-danger", lit: false, id: "supporters" },
            { title: "Orders" , icon: "shopping-cart text-success", lit: false, id: "orders" },
            { title: "Earnings", icon: "credit-card", lit: false, id: "earnings" },
            { title: "Views and engagement", icon: "line-chart text-info", lit: false, id: "engagement" },
            { title: "Notifications", icon: "envelope", lit: false, id: "notifications" },
            { title: "Delete", icon: "warning text-warning", lit: false, id: "delete" },
        ]

        v1.campaign.get(id, false).then(response => {
            v1.user.current(true).then(user => {
                this.setState({ campaign: response.campaign, user, stats: response.stats, supporters: response.supporters, files: response.files, loading: false })
            })
        }).catch(_ => {
            v1.campaign.get(id, true).then(response => {
                this.setState({ campaign: response.campaign, stats: response.stats, supporters: response.supporters, files: response.files, loading: false })
            })
        })
    }

    render() {
        let editCampaign = <></>
        switch(this.state.campaign.type) {
            case "embed":
                editCampaign = <EditEmbedCampaign user={this.state.user} campaign={this.state.campaign} reloadFn={this.reloadFn} />
            break
            case "audio":
            case "video":
            case "photobook":
                editCampaign = <EditDigitalCampaign user={this.state.user} campaign={this.state.campaign} reloadFn={this.reloadFn} />
            break
            case "service":
                editCampaign = <EditServiceCampaign user={this.state.user} campaign={this.state.campaign} reloadFn={this.reloadFn} />
            break
            default:
                editCampaign = <h1>Loading...</h1>
        }

        return this.state.loading && this.state.user._id ? <h1>Loading...</h1> : (
            <>
                
                <div class="mr-auto">
                    <div class="d-inline-block align-items-center">
                        <nav>
                            <ul class="breadcrumb fa-2x">
                                <li class="breadcrumb-item"><Link to="/creator/content"><i class="fa fa-home"></i></Link></li>
                                <li class="breadcrumb-item" aria-current="page"><Link onClick={this.reloadSettings} to={"/creator/content/"+ this.state.campaign._id}>{this.state.campaign.title}</Link></li>
                                <li class="breadcrumb-item active" aria-current="page">
                                    <Switch>
                                    <Route path="/creator/content/:id/edit">Edit content</Route>
                                    <Route path="/creator/content/:id/orders"> Orders </Route>
                                    <Route path="/creator/content/:id/supporters"> Supporters</Route>
                                    <Route path="/creator/content/:id/earning" > Earnings</Route>
                                    <Route path="/creator/content/:id/engagement"> Views and Engagement</Route>
                                    <Route path="/creator/content/:id/notifications"> Notifications</Route>
                                    <Route path="/creator/content/:id/delete"> <i class="text-danger fa fa-warning"></i> Delete</Route>
                                    <Route path="/creator/content/:id">{this.state.campaign.title}</Route>

                                    </Switch>
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>
                <div class="tab-content tabcontent-border">
                    <div class="tab-pane active" id="activeTab" role="tabpanel">
                        <div class="box-body">
                            <div class="padding-bottom-10">
                                <CampaignSummary user={this.state.user} files={this.state.files} campaign={this.state.campaign} reloadFn={this.reloadFn} />
                            </div>
                            <Switch >
                                <Route path="/creator/content/:id/edit">
                                    {editCampaign}
                                </Route>
                                <Route path="/creator/content/:id/orders">
                                    <CampaignOrders campaign={this.state.campaign} />
                                </Route>
                                <Route path="/creator/content/:id/supporters">
                                    <CampaignSupporters campaign={this.state.campaign} />
                                </Route>
                                <Route path="/creator/content/:id/earning" >
                                    <CampaignRevenueSummary campaign={this.state.campaign} />
                                </Route>
                                <Route path="/creator/content/:id/engagement">
                                    <CampaignEngagement campaign={this.state.campaign} />
                                </Route>
                                <Route path="/creator/content/:id/notifications">
                                    <CampaignNotifications campaign={this.state.campaign} reloadFn={this.reloadFn} />
                                </Route>
                                <Route path="/creator/content/:id/delete"> 
                                    <CampaignDelete campaign={this.state.campaign} />
                                </Route>
                                <Route path="/creator/content/:id">
                                    <EditCampaignLandingPage pagesFn={this.pages} user={this.state.user} campaign={this.state.campaign} /> 
                                </Route>
                            </Switch>
                        </div>
                    </div>
                </div>
        </>
       )
    }

}

export default EditCampaign;