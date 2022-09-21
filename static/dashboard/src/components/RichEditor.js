import React, { Component } from 'react';
import Editor, { createEditorStateWithText } from '@draft-js-plugins/editor';
import createEmojiPlugin from '@draft-js-plugins/emoji';
import {
    convertToRaw,
} from 'draft-js';

import editorStyles from './editorStyles.module.css';
import '@draft-js-plugins/emoji/lib/plugin.css'



class RichEditor extends Component {
    constructor(props) {
        super(props)
        this.state = {
            editorState: createEditorStateWithText(this.props.text ? this.props.text : ''),
        };

        this.onChange = this.onChange.bind(this)
        this.focus = this.focus.bind(this)

        this.emojiPlugin = createEmojiPlugin({
            useNativeArt: true,
        });
    }
    

    onChange = (editorState) => {
        this.setState({
            editorState,
        });
        const blocks = convertToRaw(editorState.getCurrentContent()).blocks;
        const value = blocks.map(block => (!block.text.trim() && '\n') || block.text).join('\n');
        return this.props.onChange ? this.props.onChange(value) : null
    };

    focus = () => {
        this.editor.focus();
    };

    render() {
        return (
            <div>
                <div className={editorStyles.editor} onClick={this.focus}>
                    <Editor
                        editorState={this.state.editorState}
                        onChange={this.onChange}
                        plugins={[this.emojiPlugin]}
                        ref={(element) => {
                            this.editor = element;
                        }}
                    />
                    <this.emojiPlugin.EmojiSuggestions />
                </div>
            </div>
        );
    }
}

export default RichEditor;