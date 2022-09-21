import React, { Component } from "react"

import Uppy from '@uppy/core'
import { DashboardModal } from '@uppy/react'
import XHRUpload from '@uppy/xhr-upload'

import '@uppy/core/dist/style.css'
import '@uppy/dashboard/dist/style.css'
import '@uppy/file-input/dist/style.css'
import '@uppy/image-editor/dist/style.css'

import './upload.css'


import v1 from "../api/v1";
import ImageEditor from "@uppy/image-editor"

class ActionButtonUpload extends Component {

    constructor(props) {

        super(props)
        this.state = {
            id: "uploader_"+ Date.now() + parseInt(Math.random()*100),
            uploadOpen: this.props.defaultOpen ? this.props.defaultOpen : false,
            previewURL: this.props.image ? this.props.image : "/assets/img/placeholder.png",
        }

        this.upload = new Uppy({
            allowMultipleUploads: false,
            restrictions: {
                maxNumberOfFiles: 1,
                minFileSize: (1024 * 1),
                maxFileSize: (1024 * 1024 * 1024 * 2),
                allowedFileTypes: this.props.allowedTypes,
            },
        })
            .use(XHRUpload, { endpoint: v1.config.uploadMediaEndpoint(this.props.type, this.props.purpose) })
        .use(ImageEditor, {
                id: 'ImageEditor',
                quality: 0.8,
                cropperOptions: {
                viewMode: 1,
                background: false,
                autoCropArea: 1,
                responsive: true,
                },
            })

        // it’s probably a good idea to clear the `<input>`
        // after the upload or when the file was removed
        // (see https://github.com/transloadit/uppy/issues/2640#issuecomment-731034781)
        this.upload.on('file-removed', () => {
            if (this.fileInput) {
                this.fileInput.value = null    
            }
            
        })
        
        this.upload.on('file-added', file => {
            console.log('file-added')
            const dashboard = this.upload.getPlugin('react:Dashboard')
            console.log(dashboard)
            //dashboard.toggleFileCard(true, file.id)
            try {
                setTimeout(_ => dashboard.openFileEditor(file), 250)
            } catch (e) {
                console.log('exception')
                console.log(e)
            }
            console.log('opened file editor')
        })

        this.upload.on('complete', () => {
            if (this.fileInput) {
                this.fileInput.value = null    
            }
            
        })

        this.upload.on('upload-success', (file, response) => {
            if (response.body._id) {
                let url = response.body.url
                if (url === "") {
                    url = v1.assets.imageURL(response.body._id, 240, 240) + '?reload'
                }

                this.setState({
                    previewURL: url,
                    uploadOpen: false,
                })
                this.upload.reset()
                return this.props.onUploaded ? this.props.onUploaded(response.body) : console.log(response.body)
                
            } else {
                alert('Upload failed with error\n:' + response.body.error)
            }
        })

        this.handleHiddenUpload = (event) => {
                const files = Array.from(event.target.files)
                const file = files[0]
            try {
                this.upload.addFile({
                    source: 'file input',
                    name: file.name,
                    type: file.type,
                    data: file
                })
                this.upload.upload()
            } catch (err) {
                if (err.isRestriction) {
                    alert('Restriction error:\n' + err)
                } else {
                    alert('An error occured!\n' + err)
                    // handle other errors
                    console.error(err)
                }
            }
            
        }

        this.handleOpenF = () => {
            if(this.props.noDashboard){
                this.fileInput = document.querySelector('#' + this.state.id)
                console.log(this.state.id, this.fileInput)

                this.fileInput.click()

            }
            this.setState({ uploadOpen: true })
        }

        this.handleCloseF = () => {
            this.setState({ uploadOpen: false })
            this.upload.reset()
        }

    }


    componentWillUnmount() {
        this.upload.close()
    }

    render() {
        return (
            <>
                { this.props.noDashboard ? <input type="file" onChange={this.handleHiddenUpload} accept={this.props.allowedTypes.join(', ')} hidden={true} class="d-none" id={this.state.id} /> : <DashboardModal
                    uppy={this.upload}
                    closeModalOnClickOutside
                    open={this.state.uploadOpen}
                    onRequestClose={this.handleCloseF}
                    proudlyDisplayPoweredByUppy={true}
                    disablePageScrollWhenModalOpen={false}
                    plugins={['ImageEditor']}
                /> }
                {this.props.hidePreview ? <></> : 
                    <a href="javascript:;" onClick={this.handleOpenF}  >
                        <img src={this.state.previewURL} alt="uploaded resource" style={{width: "100%"}}  class={this.props.previewClasses? this.props.previewClasses : "img-responsive"} />
                </a> }
                { this.props.hideButton ? <></> :
                <button class={this.state.id + " "+ (this.props.uploadBtnClasses ? this.props.uploadBtnClasses : "btn btn-default btn-block ")} onClick={this.handleOpenF} ><i class="fa fa-upload"></i> {this.props.uploadBtnText ? this.props.uploadBtnText : "Upload"} </button>
                }
                
            </>
        )
    }
}

export default ActionButtonUpload;