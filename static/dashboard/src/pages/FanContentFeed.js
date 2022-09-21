import React, { Component } from "react";
import { Link } from "react-router-dom";
import ReactTimeAgo from "react-time-ago/commonjs/ReactTimeAgo";
import v1 from "../api/v1";
import CreatorContentCard from "../components/CreatorContentCard";
import CreatorContentList from "../components/CreatorContentList";

import OptionsButtonGroup from "../components/OptionsButtongroup";


function FeedListItem(props) {
    return (
        <Link class="media media-single" to={"/@"+ props.creator.username + "/" + props.content.uri} >
            <img class="w-80 rounded" src={props.content.thumbnail} alt="..." />
                <div class="media-body">
                <h6>{props.content.title}</h6>
                <small class="text-fader"><ReactTimeAgo date={props.content.created_at} /></small>
                <small class="text-fader">By @{props.creator.username}</small>
                </div>
        </Link>
    )
}

function FeedListView(props) {
    return (
        <div class="media-list media-list-hover media-list-divided">
            {props.feed.map(entry => (<FeedListItem feed={entry.content} creator={entry.creator} />))}            
		</div>
    )
}

class FanContentFeed extends Component {

    constructor(props) {
        super(props)

        this.state = {
            feed: [],
            view: 'all',
            layout: 'grid',
            viewOptions: [
                {
                    value: 'all',
                    component: <><i class="fa fa-newspaper text-danger"></i> All</>,
                },
                {
                    value: 'explore',
                    component: <><i class="fa fa-magic text-danger"></i> Discover</>,
                },
                {
                    value: 'filtered',
                    component: <><i class="fa fa-search text-danger"></i> Filter</>,
                }
            ],
            filter: 'pay_per_view',
            filterOptions: [
                {
                    value: 'pay_per_view',
                    component: <><i class="fa fa-credit-card text-danger"></i> Paid for</>,
                },
                {
                    value: 'subscriptions',
                    component: <><i class="fa fa-lock text-danger"></i> Subscriptions</>,
                },
            ]

        }

        this.reloadView = (view) => {
            console.log(view)
            let lv = this.state.view
            this.setState({ view: view, lastView: lv })
            this.load(view, this.state.filter, true)
        }
        this.reloadView = this.reloadView.bind(this)

        this.reloadFilter = (value) => {
            console.log(value)
            let lf = this.state.filter
            let lfeed = this.state.feed
            this.setState({ filter: value, lastFilter: lf, lastFeed: lfeed, loading: true})
            this.load(this.state.view, value, false)          
        }
        this.reloadFilter = this.reloadFilter.bind(this)

        this.load = (view, filter, cached) => {
            console.log('req with '+ view +'/'+filter)
            v1.feed.load(view, filter, this.state.feed.length, cached).then(resp => {
                console.log('got response' + resp)
                this.setState({ loading: false, feed: resp.feed, lastFeed: resp.feed, filter: resp.filter, view: resp.view,  err: undefined })
                //if (cached) {
                //    this.load(false)   
                //}
            }).catch(err => {
                console.log('got error')
                console.log(err)
                this.setState({loading: false, feed: this.state.lastFeed, view: this.state.lastView, filter: this.state.lastFilter, err})
            })
        }
        this.load = this.load.bind(this)
        
    }

    render() {

        return (
            <>
                <div class="box">
                    <div class="box-header">
                        <div class="row">
                            <div class="col-md-7">
                                <h4 class="box-title">Your content feed</h4>
                            </div>
                            <div class="col-md-5">
                                <div class="box-controls pull-right">
                                    <OptionsButtonGroup items={this.state.viewOptions} item={this.state.view} onChange={this.reloadView} />
                                &nbsp;
                                </div>
                                { this.state.view === 'filtered' ?
                                <div class="box-controls pull-right mt-20">
                                    <OptionsButtonGroup items={this.state.filterOptions} item={this.state.filter} onChange={this.reloadFilter} />
                                &nbsp;
                                </div> : <></>}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="box">
                    <div class="box-body p-0">
                        {this.state.layout === 'list' ? <FeedListView feed={this.state.feed} /> :
                            <div class="container row justify-content-center" >
                                {this.state.feed.map((item, idx) => {
                                    return <CreatorContentCard  user={this.props.user} creator={item.creator} showUsername={true} content={item.content} />
                                })}
                            </div>
                        }
					</div>
                    <div class="text-center bt-1 border-light p-2">
                        <button class="btn btn-block btn-primary text-uppercase d-block font-size-12" >Load more content...</button>
                    </div>
				</div>
                 
            </>
        )
    }

}

export default FanContentFeed;