/**
 * Copyright 2018 Shift Devices AG
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Component } from 'react';
import { withTranslation } from 'react-i18next';
import MenuIcon from '../../assets/icons/menu.svg';
import { share } from '../../decorators/share';
import { translate, WithTranslation } from '../../decorators/translate';
import { SharedProps, store, toggle as toggleGuide } from '../guide/guide';
import { toggleSidebar } from '../sidebar/sidebar';
import * as style from './Spinner.css';

interface SpinnerProps {
    text?: string;
}

type Props = SpinnerProps & WithTranslation & SharedProps;

class Spinner extends Component<Props> {

    public componentWillMount() {
        document.addEventListener('keydown', this.handleKeyDown);
    }

    public componentWillUnmount() {
        document.removeEventListener('keydown', this.handleKeyDown);
    }

    private handleKeyDown = (e: KeyboardEvent) => {
        e.preventDefault();
        // @ts-ignore (blur exists only on HTMLElements)
        document.activeElement.blur();
    }

    public render() {
        const {
            t,
            text,
            guideExists,
        } = this.props;
        return (
            <div className={style.spinnerContainer}>
                <div className={style.togglersContainer}>
                    <div className={style.togglerContainer}>
                        <div className={style.toggler} onClick={toggleSidebar}>
                            <img src={MenuIcon} />
                        </div>
                    </div>
                    {
                        guideExists && (
                            <div className={style.guideToggler} onClick={toggleGuide}>
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                    strokeLinejoin="round">
                                    <circle cx="12" cy="12" r="10"></circle>
                                    <circle cx="12" cy="12" r="4"></circle>
                                    <line x1="4.93" y1="4.93" x2="9.17" y2="9.17"></line>
                                    <line x1="14.83" y1="14.83" x2="19.07" y2="19.07"></line>
                                    <line x1="14.83" y1="9.17" x2="19.07" y2="4.93"></line>
                                    <line x1="14.83" y1="9.17" x2="18.36" y2="5.64"></line>
                                    <line x1="4.93" y1="19.07" x2="9.17" y2="14.83"></line>
                                </svg>
                                {t('guide.toggle.open')}
                            </div>
                        )
                    }
                </div>
                {
                    text && text.split('\n').map((line, i) => (
                        <p key={`${line}-${i}`} className={style.spinnerText}>{line}</p>
                    ))
                }
                <div className={style.spinner}>
                    <div></div>
                    <div></div>
                    <div></div>
                    <div></div>
                </div>
                <div className={style.overlay}></div>
            </div>
        );
    }
}

const SharedSpinner = share<SharedProps, SpinnerProps & WithTranslation>(store)(Spinner);
const TranslatedSpinner = withTranslation()(SharedSpinner as any);
export { TranslatedSpinner as Spinner };
