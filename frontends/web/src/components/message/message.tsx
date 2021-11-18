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

import { CSSProperties, FunctionComponent } from 'react';
import * as styles from './message.css';

export interface Props {
    hidden?: boolean;
    type?: 'message' | 'success' | 'info' | 'warning' | 'error';
    style: CSSProperties;
}

export const Message:FunctionComponent<Props> = ({
    hidden,
    type = 'message',
    style = {},
    children,
}) => {
    if (hidden) {
        return null;
    }
    return (
        <div className={styles[type]} style={style}>
            {children}
        </div>
    );
}
