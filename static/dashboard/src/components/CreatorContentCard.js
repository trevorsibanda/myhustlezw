import { Component } from "react"
import { Link } from "react-router-dom"
import LockedPayPerViewModal from "./payments/LockedPayPerViewModal"
import SubscriptionPayModal from "./payments/SubscriptionPayModal"
import money from "./payments/Amount"

class CreatorContentCard extends Component {

    constructor(props) {
        super(props)

        this.state = {
            showPayPerViewModal: false,
            showSubscribeModal: false,
        }
        this.toggleModal = this.toggleModal.bind(this)
        this.toggleCloseModal = this.toggleCloseModal.bind(this)
    }

    toggleModal(evt) {
        if (this.props.content.subscription === 'pay_per_view' && !this.props.content.can_view) {
            this.setState({ showPayPerViewModal: !this.state.showPayPerViewModal })
            evt.preventDefault()
            return
        }
        if (this.props.content.subscription === 'fans' && !this.props.content.can_view) {
            this.setState({ showSubscribeModal: !this.state.showSubscribeModal })
            evt.preventDefault()
            return
        }
    }

    toggleCloseModal() {
        if (this.state.showPayPerViewModal) {
            this.setState({ showPayPerViewModal: false })
        } else {
            this.setState({ showSubscribeModal: false })
        }
    }


    render() {
        
        let contentLabel = this.props.content.type
        let contentAccess = this.props.content.subscription
        let contentLink = "/@" + this.props.creator.username + "/" + this.props.content.uri
        if( this.props.redirect ){
            contentLink = "/_r" + contentLink 
        }
        let contentIcon = <></>

        switch (this.props.content.type) {
            case 'video':
            case 'embed':
                contentLabel = "Video"
                break;
            case 'audio':
                contentLabel = "Audio"
                break;
            case 'image':
                contentLabel = "Image"
                break;
            case 'service':
                contentLabel = "Service"
                break;
            case 'other':
                contentLabel = "Download"
                break;
            default:
                contentLabel = 'Error'
                break;
        }

        let iconClasses = this.props.content.can_view ? "fa fa-unlock" : "fa fa-lock"

        switch (this.props.content.subscription) {
            case 'public':
                contentAccess = null
                if (this.props.content.type === 'service') {
                    contentAccess = money.formatUSD(this.props.content.price)
                }
                break;
            case 'pay_per_view':
                contentAccess = money.formatUSD(this.props.content.price)
                contentIcon = <i class={iconClasses}></i>
                break;
            case 'fans':
                contentAccess = "Fans Only"
                contentIcon = <i class={iconClasses}></i>
                break;
            default:
                contentAccess = 'Error'
                break;
        }
        
        let component = <></>
        let root = (
            <div class={"col-md-6 col-sm-6 margin-bottom-10 "+ (this.props.classes ? this.props.classes : '')}>
                <div class="media-grid">
                    <div class="media-image">
                        <Link to={contentLink} onClick={this.toggleModal}>
                            <img class="pic-1" alt="no description" src={this.props.content.preview_url} />
                        </Link>
                        <span class="media-new-label">{contentLabel}</span>
                        {contentAccess ? <span class="media-discount-label">{contentAccess}</span> : <></>}
                        {this.props.showUsername ? <span class="media-discount-label bg-success mt-20">@{this.props.creator.username}</span> : <></>}
                    </div>
                    <div class="media-content">
                        <h3 class="title"><Link onClick={this.toggleModal} to={contentLink}>{contentIcon} {this.props.content.title}</Link></h3>
                    </div>
                </div>
            </div>
        )
        
        if (this.state.showPayPerViewModal && this.props.content.subscription === 'pay_per_view' && !this.props.content.can_view) {
            component = <LockedPayPerViewModal content={this.props.content} target={contentLink} onClose={this.toggleCloseModal} user={this.props.user} creator={this.props.creator} />
        } else if (this.state.showSubscribeModal && this.props.content.subscription === 'fans' && !this.props.content.can_view) {
            component = <SubscriptionPayModal onClose={this.toggleCloseModal} target={contentLink} user={this.props.user} creator={this.props.creator} />
        }

        return <>{root}{component}</>;
    }
}

export default CreatorContentCard;