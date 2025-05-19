/**
 *
 * (c) Copyright Ascensio System SIA 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package settings

import "errors"

var (
	ErrAddressRequired = errors.New("address is required in non-demo mode")
	ErrSecretRequired  = errors.New("secret is required in demo mode")
	ErrHeaderRequired  = errors.New("header is required in demo mode")
	ErrInvalidURL      = errors.New("address must be a valid URL")
	ErrInvalidProtocol = errors.New("address must use http or https protocol")
	ErrTrailingSlash   = errors.New("address must not have a trailing slash")
	ErrHeaderTooLong   = errors.New("header must be at most 255 characters")
	ErrSecretTooLong   = errors.New("secret must be at most 255 characters")
)
