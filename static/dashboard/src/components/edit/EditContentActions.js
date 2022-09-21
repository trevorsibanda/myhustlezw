import { Component } from "react";
import SweetAlert from "react-bootstrap-sweetalert";
import CreatorRecentSupporters from "../CreatorRecentSupporters";
import PageEarnings from "../PageEarnings";
import SharePageUI from "../SharePageUI";
import EditDigitalContent from "./EditDigitalContent";
import v1 from "../../api/v1"
import { Redirect } from "react-router-dom/cjs/react-router-dom.min";

class EditContentActions extends Component {

    constructor(props){
        super(props);

        this.state = {
            action: this.props.action,
            target: '',
            show: false,
        }

        if (this.props.action && this.props.action !== '') {
            this.setState({show: true})
        }

        this.onClose = () => { this.setState({show: false}) }
        
        this.handleEvent = (evt, action) => {
            this.setState({
                action: action,
                show: true,
            })
            evt.preventDefault();
        }

        this.onDelete = (evt => {
            v1.campaign.deleteCampaign(this.props.content._id).then(resp => {
                if (resp && resp.error) {
                    alert('Failed to delete with error:\n' + resp.error)
                    return
                } else {
                    alert('Successfully deleted ' + this.props.content.type + ' - '+ this.props.content.title)
                    this.setState({action: 'redirect', target: '/@'+this.props.creator.username})
                }
            })
            evt.preventDefault()
        })
    }

    render() {
        var component, title
        switch (this.state.action) {
            case 'edit':
                title = 'Edit content'
                component = <EditDigitalContent creator={this.props.creator} content={this.props.content} onClose={this.onClose} />; 
                break;
            case 'redirect':
                component = <Redirect to={this.state.target} />
                break;
            case 'share':
                title = 'Share content'
                component =  <SharePageUI/>
                break
            case 'supporters':
                title = 'Supporters'
                console.log(this.props)
                component = <CreatorRecentSupporters user={this.props.user} creator={this.props.creator} maxShowMobile={8} grandMax={8} loadMore={true} content={this.props.content} supporters={this.props.supporters}/>;
                break
            case 'earnings':
                title = 'Earnings'
                component = <PageEarnings earnings={{zwl: 0.00, usd: 0.00}} />
                break
            case 'orders':
                title = 'Orders'
                component = <p>Orders</p>
                break
            case 'engagement':
                title = 'Views and Engagement'
                component = <div class="box box-body">
                    <p>Views and engagement tracking not yet active<br />You will get notified when this feature is now available.</p>
                    <div class="row justify-content-center">
                        <div class="col-sm-12 mt-30">
                            <button class="btn btn-default btn-block" onClick={() => { this.setState({show: false}) }}>OK !</button>
                        </div>
                        
                    </div>
                    </div>
                break
            case 'delete':
                title = 'Confirm delete?'
                component = <div class="box box-body">
                    <p>Are you sure you want to delete this content?<br />This action cannot be reversed.</p>
                    <div class="row justify-content-center">
                        <div class="col-sm-12">
                            <button class="btn btn-danger btn-block" onClick={this.onDelete}>
                                <i class="fa fa-warning"></i> Delete</button>
                        </div>
                        <div class="col-sm-12 mt-30">
                            <button class="btn btn-default btn-block" onClick={() => { this.setState({show: false}) }}>CANCEL</button>
                        </div>
                        
                    </div>
                    </div>
                break
            default:
                title = 'Unknown option'
                component = <p>Unknown option</p>
                break

        }

        return (
            <div class="box">
                <div class="box-body p-0">
                    <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={title} onCancel={this.onClose}  >
                        <div style={{textAlign: 'initial'}} > {component}</div>
                    </SweetAlert>
                    <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "share")}>
                        <div class="media align-items-center">
                            <div class="status-success"><i class="fa fa-share text-success"></i></div>
                            <div class="media-body">
                                <p class="font-size-16">
                                    <div class="hover-primary"><strong>Share this page</strong></div>
                                </p>
                            </div>
                        </div>
                        <div class="media pt-0 d-none d-md-block">
                            <p>.</p>
                        </div>
                    </a>
                    <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "edit")}>
                        <div class="media align-items-center">
                            <div class="status-success"><i class="fa fa-edit"></i></div>
                            <div class="media-body">
                                <p class="font-size-16">
                                    <div class="hover-primary"><strong>Edit content</strong></div>
                                </p>
                            </div>
                        </div>
                        <div class="media pt-0 d-none d-md-block">
                            <p>.</p>
                        </div>
                    </a>
                    <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "supporters")}>
                        <div class="media align-items-center">
                            <div class="status-success"><i class="fa fa-heart text-danger"></i></div>
                            <div class="media-body">
                                <p class="font-size-16">
                                    <div class="hover-primary"><strong>Supporters</strong></div>
                                </p>
                            </div>
                        </div>
                        <div class="media pt-0 d-none d-md-block">
                            <p>.</p>
                        </div>
                    </a>
                    {this.props.content.type === "service" ?
                        <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "orders")}>
                            <div class="media align-items-center">
                                <div class="status-success"><i class="fa fa-shopping-cart text-success"></i></div>
                                <div class="media-body">
                                    <p class="font-size-16">
                                        <div class="hover-primary"><strong>Orders</strong></div>
                                    </p>
                                </div>
                            </div>
                            <div class="media pt-0 d-none d-md-block">
                                <p>.</p>
                            </div>
                        </a> : <></>}
                    {this.props.content.subscription === "pay_per_view" || this.props.content.subscription === "subscription" ?
                        <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "earnings")}>
                            <div class="media align-items-center">
                                <div class="status-success"><i class="fa fa-credit-card"></i></div>
                                <div class="media-body">
                                    <p class="font-size-16">
                                        <div class="hover-primary"><strong>Earnings</strong></div>
                                    </p>
                                </div>
                            </div>
                            <div class="media pt-0 d-none d-md-block">
                                <p>.</p>
                            </div>
                        </a> : <></>}
                    <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "engagement")}>
                        <div class="media align-items-center">
                            <div class="status-success"><i class="fa fa-line-chart text-info"></i></div>
                            <div class="media-body">
                                <p class="font-size-16">
                                    <div class="hover-primary"><strong>Views and engagement</strong></div>
                                </p>
                            </div>
                        </div>
                        <div class="media pt-0 d-none d-md-block">
                            <p>.</p>
                        </div>
                    </a>
      
                    <a class="media-list bb-1 bb-dashed border-light" href="javascript:;" onClick={evt => this.handleEvent(evt, "delete")}>
                        <div class="media align-items-center">
                            <div class="status-success"><i class="fa fa-warning text-warning"></i></div>
                            <div class="media-body">
                                <p class="font-size-16">
                                    <div class="hover-primary"><strong>Delete</strong></div>
                                </p>
                            </div>
                        </div>
                        <div class="media pt-0 d-none d-md-block">
                            <p>.</p>
                        </div>
                    </a>
                </div>
            </div>
        )
    }
}

export default EditContentActions;