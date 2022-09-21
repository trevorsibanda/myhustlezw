import v1 from "../api/v1";
import { Component } from "react";
import CreatorContentList from "../components/CreatorContentList";
import CreatorFeaturedContent from "../components/CreatorFeaturedContent";
import CreatorRecentSupporters from "../components/CreatorRecentSupporters";
import CreatorSmallSubscribe from "../components/CreatorSmallSubscribe";
import CreatorSmallSupport from "../components/CreatorSmallSupport";
import Preloader from "../components/PreLoader";
import PublicCreatorNotFound from "./PublicCreatorNotFound";
import PublicCreatorNotVerified from './PublicCreatorNotVerified'
import CreatorSubscriptionDetails from "../components/CreatorSubscriptionDetails";


class PublicCreatorAccountPage extends Component{

    constructor(props){
        super(props)

        v1.page.track()
        v1.page.set({title: '@'+ this.props.creator.username })
    }
    
    render() {
        let component = <></>
        
        if (this.props.not_found) {
            component = <PublicCreatorNotFound />
        } else if (this.props.creator && !this.props.creator.verified) {
            component = <PublicCreatorNotVerified /> 
        }
        else if (this.props.creator && this.props.creator._id) {
            component =
                <div class="padding-bottom-40 padding-top-10" >
                    <div class="row justify-content-center" >
                        <div class="col-lg-4 col-xs-12 col-md-4 order-md-1">
                            <CreatorFeaturedContent creator={this.props.creator} featured={this.props.featured} stream_url={this.props.stream_url} />
                            {this.props.subscription.support_type === 'subscribed' ? <CreatorSubscriptionDetails subscription={this.props.subscription} creator={this.props.creator} /> :
                            (<>
                                {this.props.creator.page.allow_supporters ? <CreatorSmallSupport creator={this.props.creator} user={this.props.user} /> : <></>}
                                {this.props.creator.subscriptions.active && !this.props.creator.page.allow_supporters ? <CreatorSmallSubscribe creator={this.props.creator} user={this.props.user} /> : <></>}
                            </>)}
                            <CreatorRecentSupporters grandMax={5} maxShowMobile={2} supporters={this.props.supporters} creator={this.props.creator} user={this.props.user} />
                        </div>
                        <div class="col-lg-8 col-xs-12 col-md-8 order-md-2 mb-20">
                            <CreatorContentList loadMore={true} content={this.props.campaigns} user={this.props.user} creator={this.props.creator} />
                        </div>
                        
                    </div>
                    
                </div>
        }
        
        return (this.props.loading ? <Preloader /> : component)
    }
}

export default PublicCreatorAccountPage;