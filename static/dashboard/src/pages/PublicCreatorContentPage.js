import { Component } from "react";
import ReactPlayer from "react-player";
import v1 from "../api/v1";
import money from "../components/payments/Amount"
import CreatorBottomNav from "../components/CreatorBottomNav";
import EditContentActions from "../components/edit/EditContentActions";
import ImageViewer from "../components/ImageViewer";

import Preloader from "../components/PreLoader";
import Linkify from "react-linkify/dist/components/Linkify";
import ServicePayModal from "../components/payments/ServicePayModal";
import PublicCreatorNotFound from "./PublicCreatorNotFound";
import PublicCreatorNotVerified from "./PublicCreatorNotVerified";
import CreatorFeaturedContent from "../components/CreatorFeaturedContent";
import CreatorSmallSupport from "../components/CreatorSmallSupport";
import CreatorSmallSubscribe from "../components/CreatorSmallSubscribe";
import CreatorSubscriptionDetails from "../components/CreatorSubscriptionDetails";
import CreatorContentList from "../components/CreatorContentList";
import ReactTimeAgo from "react-time-ago/commonjs/ReactTimeAgo";
import PaymentUI from "../components/payments/PaymentUI";


class ViewServiceContentPage extends Component {
    constructor(props) {
        super(props)

        this.state = {
            fullname: props.user && props.user.fullname ? props.user.fullname : '',
            email: props.user && props.user.email ? props.user.email :  '',
            phone: props.user && props.user.phone_number ? props.user.phone_number : '',
            answer: '',
            showPayModal: false,
        }

        this.onSubmit = this.onSubmit.bind(this)
        this.validateForm = this.validateForm.bind(this)
        this.onCloseModal = this.onCloseModal.bind(this)
    }

    validateForm() {
        let phone = this.state.phone
        if (phone.length < 9) {
            alert("Please enter a valid mobile phone number")
            return false
        }

        if (this.state.fullname === '' ) {
            alert('Please enter a valid name')
            return false
        }

        if (!v1.util.validateEmail(this.state.email)) {
            alert("Please enter a valid email address")
            return false
        }

        if (this.state.answer === '') {
            alert('You left the answer field empty. Please enter a value')
            return false
        }
        return true
    }

    onSubmit(evt) {
        v1.page.event('Service Pay Modal', 'Open', this.props.creator.username)
        if (!this.validateForm()) {
            v1.page.event('Service Pay Modal', 'Fail Validation', this.props.creator.username)
            return false
        }

        this.setState({
            form: {
                fullname: this.state.fullname,
                phone: this.state.phone,
                email: this.state.email,
                answer: this.state.answer,
                question: this.props.service.service.question,
                quantity: this.props.service.service.quantity_available,
            },
            showPayModal: !this.state.showPayModal,
        })
    }

    onCloseModal() {

        this.setState({ showPayModal: !this.state.showPayModal })
        v1.page.event('Service Pay Modal', 'Close', this.props.creator.username)
    }

    render() {
       
        return (
            <div class="row ">
                {this.state.showPayModal ? <ServicePayModal creator={this.props.creator} service={this.props.service} user={this.props.user} form={this.state.form} onClose={this.onCloseModal} /> : <></>}
                <div class="col-12 padding-bottom-10">
                    <h4>{this.props.service.title} <small>fulfilled by @{this.props.creator.username} </small></h4>
                    <hr class="hr-primary" />
                </div>
                <div class="col-md-12">
                        <div class="card">
                            <div class="row">
                                <aside class="col-sm-5 border-right">
                                    <ImageViewer showThumbnails={false} images={[{ original: this.props.preview.url, thumbnail: this.props.preview.thumbnail }]} />
                                </aside>
                                <aside class="col-sm-7">
                                    <article class="card-body ">
                                    <h3 class="title mb-3">{this.props.service.title}</h3>
                                        
                                        <p class="price-detail-wrap">
                                            
                                            <span>fulfilled by {this.props.creator.fullname} </span>
                                        </p>
                                        <dl class="item-property">
                                            <dt>Price</dt>
                                            <dd>
                                                <h5 class=" text-warning">
                                                    <money.USD amount={this.props.service.price} />
                                                </h5>
                                                
                                            </dd>
                                        </dl>
                                        <dl class="item-property">
                                            <dt>Description</dt>
                                            <dd>
                                            <p><Linkify
                                componentDecorator={(decoratedHref, decoratedText, key) => (
                                    <a target="blank" style={{color: 'red', fontWeight: 'bold'}} rel="noopener" target="_blank" href={decoratedHref} key={key}>
                                        {decoratedText}
                                    </a>
                                )}
                            >{this.props.service.description}</Linkify> </p>
                                            </dd>
                                        </dl>
                                        <dl class="param param-feature">
                                            <dt>Fullfillment time</dt>
                                            <dd>3 days</dd>
                                        </dl>
                                         
                                        <dl class="param param-feature">
                                            <dt>Refund policy</dt>
                                        <dd>You can open a refund dispute if the creator does
                                            not deliver the expected service.
                                            If the dispute is not resolved within 72Hours, we will process your refund.

                                            </dd>
                                        </dl>
                        
                                        <hr/>
                                        <div class="row">
                                            <div class="col-sm-12">
                                                <dl class="param param-feature">
                                                <dt>{this.props.service.service.question}</dt>
                                                    <dd>
                                                        <input type="text" placeholder="Your answer" class="form-control" value={this.state.answer} onChange={evt => this.setState({answer: evt.target.value})} />
                                                    <small>This information will be shared with @{this.props.creator.username}</small>
                                                    </dd>
                                                </dl>
                                                    <dl class="param param-feature">
                                                        <dt>* Your fullname</dt>
                                                        <dd>
                                                    <input type="text" placeholder="Your name" class="form-control" value={this.state.fullname} onChange={evt => this.setState({fullname: evt.target.value})} />
                                                        </dd>
                                                    </dl>
                                                    <dl class="param param-feature">
                                                        <dt>* Your email</dt>
                                                        <dd>
                                                            <input type="email" placeholder="Your email" class="form-control" value={this.state.email} onChange={evt => this.setState({email: evt.target.value})} />
                                                        </dd>
                                                    </dl>
                                                    <dl class="param param-feature">
                                                        <dt>* Contact phone number</dt>
                                                        <dd>
                                                            <input type="tel" placeholder="+263..." class="form-control" value={this.state.phone} onChange={evt => this.setState({phone: evt.target.value})} />
                                                        </dd>
                                                    </dl>
                                            </div>
                                            
                                        </div>
                                        <hr/>
                                        <p><strong>Prices are shown in USD for convenience.</strong></p> 
                                        
                                        <p>You can pay <span class="text-danger"><money.USD amount={this.props.service.price} /></span> 
                                        or <span class="text-danger"><money.ZWL usd={this.props.service.price} /></span></p>
                                        <p><span>Only {this.props.service.service.quantity_available} orders left</span></p>
                                    <button onClick={this.onSubmit} class="btn btn-lg btn-primary btn-block text-uppercase"> Place your order </button>
                                    </article>
                                </aside>
                            </div>
                        </div>

                </div>
            </div>
        )
    }
}

function ViewVideoContentPage(props) {
    var config = {}
    if (props.video.processed && !window.supportsHLS()) {
        v1.page.sysEvent('Video Player', 'HLS', 'forceHLS')
        config = { file: { forceHLS: true } }
    } else {
        v1.page.sysEvent('Video Player', 'HLS', 'nativeHLS')
    }
    return (
        <>
        <div class="row">
            <div class="col-md-12">
                <h4>{props.content.title}</h4>
                <hr class="hr-primary" />
            </div>
        </div>
        <div class="row justify-content-center" >
            <div class="col-md-12 col-sm-12">
                <div className="player-wrapper">
                    <ReactPlayer config={config} controls={true} width='100%' height='100%' url={props.stream_url} className="react-player" />
                </div>
            </div>
        </div>
        </>
    )
}

function ViewEmbedContentPage(props) {
    return (
        <>
        <div class="row">
            <div class="col-md-12">
                <h4>{props.embed.title}</h4>
                <hr class="hr-primary" />
            </div>
        </div>
        <div class="row justify-content-center" >
            <div class="col-md-12 col-sm-12">
                <div className="player-wrapper">
                    <ReactPlayer width='100%' controls={true} height='100%' playing={true} url={"https://youtu.be/" + props.embed.remote_id} className="react-player" />
                </div>
            </div>
        </div>
        </>
    )
}

function ViewImagesContentPage(props) {
    let processed = props.photos.map(photo => {
        return {
            original: photo.url,
            thumbnail: photo.thumbnail,
            caption: photo.caption,
        }
    })
    return (
        <>
        <div class="row">
            <div class="col-md-12">
                <h4>{props.content.title}</h4>
                <hr class="hr-primary" />
            </div>
        </div>
        <div class="row justify-content-center" >
            <div class="col-md-12 col-sm-12">
                    <ImageViewer images={processed} />
            </div>
        </div>
        </>
        
    )
}


function ContentTabPages(props) {
    return (
        <div class="row">
            <div class="col-md-12">
                <nav class="nav nav-pills bg-dark nav-fill nav-justified">
                    {
                        props.components.map(component => {
                            return <a onClick={evt => { props.onClick(component.id); evt.preventDefault(); }} class={"flex-sm-fill text-sm-center nav-link" + (component.id == props.id ? " active": "")} aria-current="page" href="javascript:;"><i class={"fa fa-" + component.icon} ></i> {component.title}</a>
                    
                        })
                    }
                </nav>
                {
                    props.components.map(component => {
                        return (
                            
                        <div class={"mt-20 " + (component.id === props.id ? "" : "d-none")} >
                                {component.component}
                        </div>)
                
                    })
                }
                
            </div>
        </div>
    )
}

function ContentDetails(props){
    return (
        <div class="row">
            <div class="col-md-12">
            <table class="table table-responsive table-striped">
            <tbody>
              <tr>
                <td>Title </td>
                <td>{props.content.title}</td>

              </tr>
              <tr>
                <td>Description</td>
                <td>{props.content.description}</td>
              </tr>
              <tr>
                <td>Created</td>
                <td><ReactTimeAgo date={props.content.created_at}/></td>
              </tr>
              <tr>
                <td>Subscription</td>
                <td>{props.content.subscription}</td>
              </tr>
              {props.content.subscription === 'pay_per_view' ?
              <tr>
                <td>Price </td>
                            <td><money.USD amount={props.content.price} /> / <money.ZWL usd={props.content.price} /> </td>
              </tr> : <></>}
              {props.content.type === 'service' ? <>
              <tr>
                <td>Instructions</td>
                <td>{props.content.service.instructions}</td>
              </tr>
              <tr>
                <td>Information required</td>
                <td>{props.content.service.question}</td>
            
              </tr>
              <tr>
                <td>Items left</td>
                <td>{props.content.service.quantity_left}  at time of page load</td>
              </tr></> : <></>}
            </tbody>
            </table>
            </div>
        </div>
    )
}


function ContentPayToView(props) {
    let pay_to_view = [ {
        original: props.content.preview_url,
        thumbnail: props.content.preview_url,
        caption: 'Pay to unlock content'
        }]
    let purpose = 'pay_per_view'
    if (props.content.subscription === 'fans') {
        purpose = 'subscribe'
    }
    return (
        <>
        <div class="row">
            <div class="col-md-12">
                <h4><i class="fa fa-lock"></i> {props.content.title}</h4>
                <h5>Pay USD<money.USD amount={props.content.price} /> / <money.ZWL usd={props.content.price} /> to unlock this content </h5>
                <hr class="hr-primary" />
                
            </div>
        </div>
        <div class="row justify-content-center" >
            <div class="col-md-12 col-sm-12">
                    <ImageViewer images={pay_to_view} />
                    <p>
                        <PaymentUI purpose={purpose} content={props.content} target={props.target} creator={props.creator} amount={props.content.price} user={props.user} />
                    </p>
            </div>
        </div>
        </>
    )
}

function ViewDownloadPage(props) {
    return (
        <>
            View download
        </>
    )
}

class PublicCreatorContentPage extends Component{
    constructor(props) {
        super(props);
        this.state = {
            loading: true,
            not_found: false,
            username: this.props.match.params.username,
            content_id: this.props.match.params.content_id,
            active_tab: 'details',
            files: [],
        }

        if(window.pregenerated_content && window.pregenerated_content.content && window.pregenerated_content.content.uri === this.state.content_id){
            this.state = { ...this.state, ...window.pregenerated_content, loading: false }
            console.log(this.state)
            window.pregenerated_content = undefined
            v1.page.set(this.state.page)
        }

        this.contentPage = this.contentPage.bind(this);
        this.switchTab = (active_tab) => {
            this.setState({active_tab})
        }

        v1.public.getCampaign(this.state.username, this.state.content_id, true).then(res => {
            this.setState(res)
            v1.page.set({title: res.page.title})
            switch(this.state.content.type) {
                case 'video': this.setState({ files: [res.video] }); break
                case 'audio': this.setState({ files: [res.audio] }); break
                case 'image': this.setState({ files: res.photos }); break
                case 'service': this.setState({ files: res.photos }); break
                default:
                    this.setState({ files: [] }); break
            }
            this.setState({
                loading: false,
                not_found: false,
            })
            window.pregenerated_content = undefined
        }).catch(err => {
            this.setState({
                loading: false,
                not_found: true,
            })
            window.pregenerated_content = undefined
            //alert('Failed to load content with error: '+ err.error)
        })

        v1.page.track()
    }
    
    contentPage(canView) {
        if (!canView) {
            return <ContentPayToView content={this.state.content} user={this.state.user} creator={this.state.creator} />
        }
        v1.page.sysEvent('Creator Content', 'View Page', this.state.content.type)
        switch(this.state.content.type) {
            case 'video':
            case 'audio':
            return <ViewVideoContentPage content={this.state.content} video={this.state.video} stream_url={this.state.stream_url} />
            case 'image':
            return <ViewImagesContentPage content={this.state.content} photos={this.state.photos} />
            case 'embed':
            return <ViewEmbedContentPage embed={this.state.content} />
            case 'service':
            return <ViewServiceContentPage user={this.state.user} service={this.state.content} creator={this.state.creator} preview={this.state.preview}  />
            case 'other':
            return <ViewDownloadPage files={this.state.files} content={this.state.content} />
            default:
                return <div>Unknown content type {this.state.content.type}</div>
        }
    }


    render() {
        let tabComponents = [
            {
                component: <ContentDetails creator={this.state.creator} content={this.state.content} />,
                id: 'details',
                icon: 'info',
                title: 'Details',
            },
            {
                title: 'Similar content',
                component: <CreatorContentList loadMore={false} content={this.state.recommendations} redirect={true} user={this.state.user} creator={this.state.creator} />,
                icon: 'globe',
                id: 'discover',
            }
        ]

        if (this.state.creator && (this.state.creator._id === this.state.user._id)) {
            tabComponents.push({
                title: 'Edit ',
                icon: 'pencil',
                id: 'edit',
                component: <EditContentActions user={this.state.user} supporters={this.state.supporters} content={this.state.content} creator={this.state.creator} files={this.state.files} />,
            })
            tabComponents.reverse()
        }

        let component = <></>
        
        if (this.state.not_found) {
            component = <PublicCreatorNotFound />
        } else if (this.state.creator && !this.state.creator.verified) {
            component = <PublicCreatorNotVerified /> 
        } else if (this.state.creator && this.state.creator._id) {
            let canView = (this.state.content.can_view || this.state.content.allowed )&& this.state.content.subscription !== 'public'
            if (this.state.content.type === 'service') {
                canView = true
            }
            component = (
            <>
                <div class="padding-bottom-40 padding-top-10" >
                        <div class="row justify-content-center" >
                            <div class="col-md-4 d-none d-md-block">
                                <CreatorFeaturedContent creator={this.state.creator} featured={this.state.featured} />
                                {this.state.subscription.support_type === 'subscribed' ? <CreatorSubscriptionDetails subscription={this.state.subscription} creator={this.state.creator} /> :
                                (<>
                                    {this.state.creator.page.allow_supporters ? <CreatorSmallSupport creator={this.state.creator} user={this.state.user} /> : <></>}
                                    {this.state.creator.subscriptions.active && !this.state.creator.page.allow_supporters ? <CreatorSmallSubscribe creator={this.state.creator} user={this.state.user} /> : <></>}
                                    </>)}
                            </div>    
                            <div class="col-md-8">
                                {this.contentPage(canView)}
                                {canView ? <ContentTabPages
                                    id={this.state.active_tab}
                                    components={tabComponents}
                                    onClick={this.switchTab}
                                /> : <></> }
                            </div> 
                            
                    </div>
                    
                </div>
                <CreatorBottomNav user={this.state.user} creator={this.state.creator} content={this.state.content} />
            </>)
        }
        return this.state.loading ? <Preloader /> : component
    }
}

export default PublicCreatorContentPage;