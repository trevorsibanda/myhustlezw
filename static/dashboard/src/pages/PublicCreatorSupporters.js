import { Component } from "react";
import v1 from "../api/v1";
import CreatorBottomNav from "../components/CreatorBottomNav";
import CreatorPage from "./CreatorPage";
import CreatorRecentSupporters from "../components/CreatorRecentSupporters";

import Preloader from "../components/PreLoader";


class PublicCreatorSupporters extends Component{
    constructor(props) {
        super(props);
        this.state = {
            loading: false,
            next: 1,
        }

        v1.page.set({title: 'Public supporters of @'+ this.props.creator.username})
        v1.page.track()

        this.loadMoreContent = this.loadMoreContent.bind(this)
    }

    loadMoreContent() {
        alert('todo: load more content')
    }
    
    render() {
        return this.state.loading ? <Preloader /> : (

            <>
                <div class="padding-bottom-40 padding-top-10 container" >
                    <div class="row justify-content-center" >
                        <div class="col-lg-10 col-xs-12 col-md-10 order-md-1">
                            <CreatorRecentSupporters maxShowMobile={100} loadMore={false} supporters={this.props.supporters} user={this.props.user} creator={this.props.creator} />
                        </div>
                        
                    </div>
                    <div class="row justify-content-center" >
                        <div class="col-lg-10 col-xs-12 col-md-10 order-md-1">
                            <div class="text-center">
                                <button onClick={this.loadMoreContent} class="btn btn-primary"><i class="fa fa-history"></i> Load More ...</button>
                            </div>
                        </div>
                    </div>
                    
                </div>
                <CreatorBottomNav user={this.props.user} creator={this.props.creator} />
            </>
        )
    }
}

export default PublicCreatorSupporters;