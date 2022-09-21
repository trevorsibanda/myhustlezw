import { Component } from "react";
import { Link } from "react-router-dom";
import v1 from "../api/v1";
import CreatorBottomNav from "../components/CreatorBottomNav";
import CreatorContentList from "../components/CreatorContentList";
import CreatorPage from "./CreatorPage";

import Preloader from "../components/PreLoader";


class PublicCreatorContentList extends Component{
    constructor(props) {
        super(props);
        this.state = {
            loading: true,
            username: this.props.match.params.username,
            content_id: this.props.match.params.content_id,
        }

        v1.page.set({title: 'Content by @'+this.state.username})

        v1.public.getCreator(this.state.username).then(res => {
            if (res.error) {
                this.setState({
                    loading: false,
                    error: res.error
                })
                alert('Failed to load with error '+ res.error)
            }
            console.log(res)
            this.setState(res)
            this.setState({loading: false})
        }).catch(err => {
            alert('Failed to load with error ' + err)
            return
        })

        v1.page.track()
        this.loadMoreContent = this.loadMoreContent.bind(this)
    }

    loadMoreContent() {
        v1.page.event('Creator  Content', 'Load More Button Click', this.state.creator.username)
        alert('todo: load more content')
    }
    
    render() {
        return this.state.loading ? <Preloader /> : (
   <><div class="padding-bottom-40 padding-top-10 container" >
                    <h4>@{this.state.creator.username}'s page</h4>
                    <hr class="hr-primary" />
                    <div class="row justify-content-center" >
                        <div class="col-lg-10 col-xs-12 col-md-10 order-md-1">
                            <CreatorContentList loadMore={false} content={this.state.campaigns} user={this.state.user} creator={this.state.creator} />
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
                <CreatorBottomNav user={this.state.user} creator={this.state.creator} />
            </>
        )
    }
}

export default PublicCreatorContentList;